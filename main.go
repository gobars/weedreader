package main

import (
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomarks/ruyi/pkg/file"
	"log"
	"os"
	"time"
)

func main() {
	//dir := "/dir_李标测试/子目录"
	//fileName := "测试机密的图片_rename.png"

	//dir := "/dir_李标测试/子目录2"
	//fileName := "Nginx"

	dir := "/dir_李标测试/子目录3"
	fileName := "MSSM-Auth-server.zip"

	entry, finalData, err := ReadFile(dir, fileName)
	if err != nil {
		log.Fatal(err)
	}
	newFileName := time.Now().String() + entry.FullPath.Name()
	writeToFile(newFileName, finalData)

	// 3,13bfdd152e.png
	fid, _ := needle.ParseFileIdFromString("3,13bfdd152e.png")
	n, err := Needle(fid)
	writeToFile(fid.String(), n.Data)
}

func ReadFile(dir string, fileName string) (*filer.Entry, []byte, error) {
	// Read file meta information from MySQL
	entry, err := FindDataEntry(dir, fileName)
	if err != nil {
		log.Printf("FindDataEntry error: %v", err)
		return entry, nil, err
	}
	// According to file meta, read files from idx and read files from dat files
	finalData, err := ReadFileData(entry)
	if err != nil {
		log.Printf("ReadFileData error: %v", err)
		return entry, nil, err
	}
	return entry, finalData, nil
}

// ReadFileData According to file meta, read files from idx and read files from dat files
func ReadFileData(entry *filer.Entry) ([]byte, error) {
	var finalData []byte
	for _, chunk := range entry.Chunks {
		parts, err := DataFromVolumeFile(chunk)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, b := range parts {
			finalData = append(finalData, b)
		}
	}
	return finalData, nil
}

// FindDataEntry Read file meta information from MySQL
func FindDataEntry(dir string, fileName string) (*filer.Entry, error) {
	fm, err := FindFileMeta(dir, fileName)
	if err != nil {
		return nil, err
	}
	if fm.Meta == nil {
		return nil, fmt.Errorf("file not exist dir:%s, file name:%s", dir, fileName)
	}
	entry := &filer.Entry{
		FullPath: util.NewFullPath(dir, fileName),
	}
	err = entry.DecodeAttributesAndChunks(util.MaybeDecompressData(fm.Meta))
	if err != nil {
		return nil, err
	}
	log.Printf("数据为：%+v", entry)
	return entry, nil
}

func getMaybeDecryptData(chunk *filer_pb.FileChunk, encryptData []byte) ([]byte, error) {
	if chunk.GetCipherKey() == nil {
		return encryptData, nil
	}
	// 加密的数据需要解密
	return util.Decrypt(encryptData, chunk.GetCipherKey())
}

func writeToFile(newFileName string, finalData []byte) {
	file.MakeSureFile(newFileName)
	os.WriteFile(newFileName, finalData, os.ModePerm)
}
