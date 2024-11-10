package main

import (
	"github.com/hdt3213/rdb/parser"
	"os"
	"strings"
	"time"
)

/*
* RDB 文件格式概述
* 以下是 RDB 文件的各个部分，按顺序排列：
* 	标题部分
* 	元数据部分
* 	数据库部分
* 	文件结束部分
 */

const (
	opCodeModuleAux    byte = 247 /* Module auxiliary data. */
	opCodeIdle         byte = 248 /* LRU idle time. */
	opCodeFreq         byte = 249 /* LFU frequency. */
	opCodeAux          byte = 250 /* RDB aux field. */
	opCodeResizeDB     byte = 251 /* Hash table resize hint. */
	opCodeExpireTimeMs byte = 252 /* Expire time in milliseconds. */
	opCodeExpireTime   byte = 253 /* Old expire time in seconds. */
	opCodeSelectDB     byte = 254 /* DB number of the following keys. */
	opCodeEOF          byte = 255
)

const (
	typeString = iota
	typeList
	typeSet
	typeZset
	typeHash
	typeZset2 /* ZSET version 2 with doubles stored in binary. */
	typeModule
	typeModule2 // Module value parser should be registered with Decoder.WithSpecialType
	_
	typeHashZipMap
	typeListZipList
	typeSetIntSet
	typeZsetZipList
	typeHashZipList
	typeListQuickList
	typeStreamListPacks
	typeHashListPack
	typeZsetListPack
	typeListQuickList2
	typeStreamListPacks2
	typeSetListPack
)

func loadRdbFileIntoKVMemoryStore() {
	//parseDB()

	// todo-w 自己实现 rdb 文件的读取逻辑
	dec := NewDecoder(openRDBFile())
	data, err := dec.parseRDB()
	if err != nil && err.Error() != "EOF" {
		logger.Error(err.Error())
		return
	}

	for i := 0; i < len(data); i += 2 {
		key := data[i]
		value := data[i+1]
		SETsMu.Lock()
		SETs[key] = &Entry{
			Value:       value,
			TimeCreated: time.Now(),
			ExpiryInMS:  time.Time{},
		}
		SETsMu.Unlock()
	}
}

func openRDBFile() *os.File {
	dir, exist := Configs["dir"]
	if !exist {
		logger.Error("Can not found db directory")
		dir = "./"
	}

	fileName, exist := Configs["dbfilename"]
	if !exist {
		logger.Error("Can not found db file name")
		fileName = "dumb.rdb"
	}

	filePath := strings.Join([]string{dir, fileName}, "/")

	rdbFile, err := os.Open(filePath)
	if err != nil {
		logger.Error("Can not found %s\n", filePath)
	}
	defer func() {
		_ = rdbFile.Close()
	}()
	return rdbFile
}

func parseDB() {
	dir, exist := Configs["dir"]
	if !exist {
		logger.Error("Can not found db directory")
		dir = "./"
	}

	fileName, exist := Configs["dbfilename"]
	if !exist {
		logger.Debug("Can not found db file name")
		fileName = "dumb.rdb"
	}

	filePath := strings.Join([]string{dir, fileName}, "/")

	rdbFile, err := os.Open(filePath)
	if err != nil {
		logger.Error("Can not found %s\n", filePath)
		return
	}
	defer func() {
		_ = rdbFile.Close()
	}()
	decoder := parser.NewDecoder(rdbFile)
	err = decoder.Parse(func(o parser.RedisObject) bool {
		switch o.GetType() {
		case parser.StringType:
			str := o.(*parser.StringObject)
			if expire := str.Expiration; expire == nil {
				SETs[str.Key] = &Entry{
					Value:       string(str.Value),
					ExpiryInMS:  time.Time{},
					TimeCreated: time.Now(),
				}
			} else {
				SETs[str.Key] = &Entry{
					Value:       string(str.Value),
					ExpiryInMS:  *expire,
					TimeCreated: time.Now(),
				}
			}
		}
		// return true to continue, return false to stop the iteration
		return true
	})
	if err != nil {
		logger.Error(err.Error())
	}

}
