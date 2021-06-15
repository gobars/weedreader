package main

import (
	"database/sql"
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

func FindFileMeta(dir string, fileName string) (FileMeta, error) {
	var fm FileMeta
	err := sqx.NewSQL("select * from filemeta where dirhash = ? and name = ? and directory = ?", util.HashStringToLong(dir), fileName, dir).QueryAsBeans(db, &fm)
	if err != nil {
		return FileMeta{}, err
	}
	return fm, nil
}
