package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/bingoohuang/gg/pkg/ctl"
	"github.com/bingoohuang/golog"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/util"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const projectName = "weedreader"
const defaultConfFile = projectName + ".yml"

var CONFIG *Config

//go:embed initassets
var initAssets embed.FS

func main() {
	prepare()
	// MySQL 数据源
	InitDataSource()

	// 读 filer 加密文件案例
	filerDemo()
	// 读 fid 文件案例
	fidDemo()
}

func fidDemo() {
	// 读 fid 文件案例
	for _, fidStr := range []string{"3,13bfdd152e", "4,1406cde3d0"} {
		fid, _ := needle.ParseFileIdFromString(fidStr)
		n, _ := Needle(fid)
		log.Println(mime.TypeByExtension(strings.ToLower(filepath.Ext(string(n.Name)))))
		writeToFile("./tmp/"+time.Now().String()+fid.String()+string(n.Name), n.Data)
	}
}

func filerDemo() {

	arr := []util.FullPath{
		util.FullPath("/buckets/b_test/1.png"),
		//util.NewFullPath("/", "1.png"),
		//util.NewFullPath("/dir_李标测试/子目录", "测试机密的图片_rename.png"),
		//util.NewFullPath("/dir_李标测试/子目录2", "Nginx"),
		//util.NewFullPath("/dir_李标测试/子目录3", "MSSM-Auth-server.zip")
	}

	for _, f := range arr {
		entry, finalData, err := ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		writeToFile("./tmp/"+time.Now().String()+entry.FullPath.Name(), finalData)
	}

}

func prepare() {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	confFile := f.String("c", "", "Filename of configuration in yaml format, default to "+defaultConfFile)
	port := f.String("p", "", "TCP ports, comma separated for multiple")
	initing := f.Bool("init", false, "init sample weedreader.yml and then exit")
	version := f.Bool("v", false, "show version info and exit")
	_ = f.Parse(os.Args[1:]) // Ignore errors; f is set for ExitOnError.
	ctl.Config{
		Initing:      *initing,
		PrintVersion: *version,
		VersionInfo:  projectName + "v0.0.1 init",
		InitFiles:    initAssets,
	}.ProcessInit()
	CONFIG = ParseConfFile(*confFile, *port)
	golog.SetupLogrus(golog.Spec(fmt.Sprintf("file=.%vlogs%v%s.log", string(os.PathSeparator), string(os.PathSeparator), projectName)))
}
