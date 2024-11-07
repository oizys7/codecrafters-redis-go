package main

import (
	"fmt"
	"os"
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

func loadRdbFileIntoKVMemoryStore() {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", *dir, *dbFileName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(content) == 0 {
		return
	}

	line := parseTable(content)
	key := line[4 : 4+line[3]]
	value := line[5+line[3]:]

	SETsMu.Lock()
	SETs[string(key)] = &Entry{
		Value:       string(value),
		TimeCreated: time.Now(),
		ExpiryInMS:  time.Time{},
	}
	defer SETsMu.Unlock()
}

func sliceIndex(data []byte, sep byte) int {
	for i, b := range data {
		if b == sep {
			return i
		}
	}
	return -1
}
func parseTable(bytes []byte) []byte {
	start := sliceIndex(bytes, opCodeResizeDB)
	end := sliceIndex(bytes, opCodeEOF)
	return bytes[start+1 : end]
}
func readFile(path string) string {
	c, _ := os.ReadFile(path)
	key := parseTable(c)
	str := key[4 : 4+key[3]]
	return string(str)
}
