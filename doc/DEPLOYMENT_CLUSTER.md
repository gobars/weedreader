# 集群部署（单机多进程方式）
## 架构图
![img.png](img_architecture.png)
## 部署图
![img_1.png](img_deployment.png)
## 部署步骤
### 规划部署机器磁盘大小、磁盘地址、服务端口

- 需要部署三类服务

|服务|个数|备注|
| ---| ---| ---|
|maser server|3|如：9331、9332、9333|
|volume server|3|如：9431、9432、9433|
|filer server|3|如：9531、9532、9533|

### 部署 seaweedfs 应用
#### 准备目录
> 为服务提供数据和日志存储

|服务|用途|备注|
| ---| ---| ---|
|maser server|data|每个进程需要单独的目录|
|maser server|log|每个进程需要单独的目录|
|volume server|data|每个进程需要单独的目录|
|volume server|log|每个进程需要单独的目录|
|filer server|log|每个进程需要单独的目录|

>示例
```shell 
mkdir -p data/m1 data/m2 data/m3
mkdir -p data/v1 data/v2 data/v3
mkdir -p log/m1 log/m2 log/m3
mkdir -p log/v1 log/v2 log/v3
mkdir -p log/f1 log/f2 log/f3
```


#### 准备应用程序
从 [seaweedfs 下载地址](https://github.com/chrislusf/seaweedfs/releases) 下载适合的版本

#### 部署 master server
##### 注意事项
- `logdir` 是 `weed` 的参数，应该紧邻 `weed` 之后
- 如果在同一台机器部署多个进程，注意 `logdir` 区分开，不要放在一起，避免 `soft link` 覆盖
- 如果机器存在多网卡、多`IP`，需要指定 `ip` 参数，避免因 `ip` 识别错误而无法组成集群
> 部署命令
```shell
weed -logdir={日志目录} master -mdir={数据目录} -ip={本机IP} -port={端口} -defaultReplication=001 -peers={组成集群的IP:端口(多个用逗号分隔)}
```
部署完成后可以访问 http://ip:port/ 查看集群 master server
> 示例
```shell
export WEED_HOME=~/Downloads/weed
nohup weed -logdir=$WEED_HOME/log/m1 master -mdir=$WEED_HOME/data/m1 -ip=localhost -port=9331 -defaultReplication=001 -peers=localhost:9331,localhost:9332,localhost:9333   > $WEED_HOME/log/m1/master.out &
nohup weed -logdir=$WEED_HOME/log/m2 master -mdir=$WEED_HOME/data/m2 -ip=localhost -port=9332 -defaultReplication=001 -peers=localhost:9331,localhost:9332,localhost:9333   > $WEED_HOME/log/m2/master.out &
nohup weed -logdir=$WEED_HOME/log/m3 master -mdir=$WEED_HOME/data/m3 -ip=localhost -port=9333 -defaultReplication=001 -peers=localhost:9331,localhost:9332,localhost:9333   > $WEED_HOME/log/m3/master.out &
```
[访问 master server](http://localhost:9333/)

#### 部署 master server
##### 注意事项
同上
> 部署命令
```shell
weed -logdir={日志目录} volume -mdir={数据目录} -ip={本机IP} -port={端口} -max={volume最大个(默认7个)} -mserver={master集群的IP:端口(多个用逗号分隔)}
```
部署完成后可以访问 http://ip:port/ 查看集群 master server
> 示例
```shell
export WEED_HOME=~/Downloads/weed
nohup weed -logdir=$WEED_HOME/log/v1 volume -dir=$WEED_HOME/data/v1/ -ip=localhost -port=9431 -max=20 -mserver=localhost:9331,localhost:9332,localhost:9333 > $WEED_HOME/log/v1/volume.out &
nohup weed -logdir=$WEED_HOME/log/v2 volume -dir=$WEED_HOME/data/v2/ -ip=localhost -port=9432 -max=20 -mserver=localhost:9331,localhost:9332,localhost:9333 > $WEED_HOME/log/v2/volume.out &
nohup weed -logdir=$WEED_HOME/log/v3 volume -dir=$WEED_HOME/data/v3/ -ip=localhost -port=9433 -max=20 -mserver=localhost:9331,localhost:9332,localhost:9333 > $WEED_HOME/log/v3/volume.out &
```
[访问 volume server](http://localhost:9433/ui/index.html)
#### 部署 filer server
##### 注意事项
- 为保证自身数据安全开启数据加密，参数 `-encryptVolumeData`
##### 准备 `MySQL`
`filer` 使用 `MySQL` 存储对象元数据信息，需要准备一个可以访问的数据库以及用户
- 数据库
- user
- 表 (不存在会自动创建)
> 示例
```sql
create database weed_cluster_0 character set utf8mb4;
create user 'weed_cluster_0' identified by 'weed_cluster_0';
grant all privileges on weed_cluster_0.* to weed_cluster_0;
 
CREATE TABLE IF NOT EXISTS weed_cluster_0.filemeta (
  dirhash     BIGINT         COMMENT 'first 64 bits of MD5 hash value of directory field',
  name        VARCHAR(1000)  COMMENT 'directory or file name',
  directory   TEXT           COMMENT 'full path to parent directory',
  meta        LONGBLOB,
  PRIMARY KEY (dirhash, name)
) DEFAULT CHARSET=utf8;
```
##### 准备 `filer server` 配置文件
|事项|说明|
|---|---|
|文件名|`filer.toml`|
|配置文件初始化方法|`weeb scaffold -config filer -output="."`|
|配置文件生效位置和顺序|The configuration file "filer.toml" is read from ".", "$HOME/.seaweedfs/", "/usr/local/etc/seaweedfs/", or "/etc/seaweedfs/", in that order.|

###### 初始化配置文件
```shell
weeb scaffold -config filer -output="."
```
> 结果
```shell
-rw-r--r--   1 anan  staff  5697  6 16 10:42 filer.toml
```
###### 配置
关闭默认的`[leveldb2]`配置，开启 `[mysql]` 配置
>示例
> [filer.toml](filer.toml)
###### 应用配置文件
配置文件被应用的顺序为：
1. "."
2. "$HOME/.seaweedfs/", 
3. "/usr/local/etc/seaweedfs/", 
4. "/etc/seaweedfs/"

建议，将配置好的配置文件移动到应用程序当前目录或者，"$HOME/.seaweedfs/" 目录下
###### 启动 `filer server`
> 部署命令
```shell
weed -logdir={日志目录} filer -encryptVolumeData -ip={本机IP} -port={端口} -mserver={master集群的IP:端口(多个用逗号分隔)} -peers={组成集群的IP:端口(多个用逗号分隔)}
```
部署完成后可以访问 http://ip:port/ 查看集群 filer server
> 示例
```shell
export WEED_HOME=~/Downloads/weed
nohup weed -logdir=$WEED_HOME/log/f1 filer -encryptVolumeData -ip=localhost -port 9531 -master=localhost:9331,localhost:9332,localhost:9333 -peers=localhost:9531,localhost:9532,localhost:9533  > $WEED_HOME/log/f1/filer.out &
nohup weed -logdir=~$WEED_HOME/log/f2 filer -encryptVolumeData -ip=localhost -port 9532 -master=localhost:9331,localhost:9332,localhost:9333 -peers=localhost:9531,localhost:9532,localhost:9533  > $WEED_HOME/log/f2/filer.out &
nohup weed -logdir=~$WEED_HOME/log/f3 filer -encryptVolumeData -ip=localhost -port 9533 -master=localhost:9331,localhost:9332,localhost:9333 -peers=localhost:9531,localhost:9532,localhost:9533  > $WEED_HOME/log/f3/filer.out &
```
[访问 filer server](http://localhost:9533/)