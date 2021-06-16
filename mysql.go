package main

import (
	"database/sql"
	"fmt"
	"github.com/bingoohuang/gg/pkg/sqx"
	"github.com/chrislusf/seaweedfs/weed/util"
	"log"
)

type FileMeta struct {
	DirHash   int64
	Name      string
	Directory string
	Meta      []byte
}

var db *sql.DB

func InitDataSource() {
	d, err := sql.Open("mysql", CONFIG.MysqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	db = d
}

func FindFileMeta(bucketName, dir, fileName string) (FileMeta, error) {
	tableName := "filemeta"
	if bucketName != "" {
		tableName = bucketName
		log.Printf("FindDefaultFileMeta from table %s", tableName)
	}
	var fm FileMeta
	sql := fmt.Sprintf("select * from %s where dirhash = ? and name = ? and directory = ?", bucketName)
	err := sqx.NewSQL(sql, util.HashStringToLong(dir), fileName, dir).QueryAsBeans(db, &fm)
	if err != nil {
		return FileMeta{}, err
	}
	return fm, nil
}
