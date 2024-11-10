package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/logging"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var port = flag.String("port", "6379", "port to listen on")
var dir = flag.String("dir", ".", "Directory to store RDB file")
var dbFileName = flag.String("dbfilename", "dump.rdb", "RDB file name")
var replicaof = flag.String("replicaof", "", "replicaof host:port")

// var logLevelStr = flag.String("loglevel", "INFO", "log print level")
var logLevel = flag.Int64("loglevel", 1, "log print level: 0 debug 1 info 2 warning 3 error 4 fatal 5 off")
var logger = logging.Logger{}

var Configs = map[string]string{}
var ConfigsMu = sync.RWMutex{}

var InfoSet = map[string]map[string]string{}
var InfoSetMu = sync.RWMutex{}

func initConfigs() {
	// 解析命令行参数
	flag.Parse()

	logger = *logging.New(int(*logLevel))
	ConfigsMu.Lock()
	Configs["loglevel"] = strconv.FormatInt(*logLevel, 10)

	Configs["port"] = *port
	Configs["dir"] = *dir
	Configs["dbfilename"] = *dbFileName
	ConfigsMu.Unlock()

	InfoSetMu.Lock()
	if *replicaof == "" {
		InfoSet["REPLICATION"] = map[string]string{"role": "master"}
	} else {
		InfoSet["REPLICATION"] = map[string]string{"role": "slave"}
		masterAndSlaveReplication()
	}

	// 初始化随机数生成器
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成 40 个字符的伪随机字母数字字符串
	randomString := GenerateRandomString(40)
	InfoSet["REPLICATION"]["masterReplId"] = randomString
	InfoSet["REPLICATION"]["masterReplOffset"] = "0"

	InfoSetMu.Unlock()
}

// GenerateRandomString 生成指定长度的伪随机字母数字字符串
func GenerateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func masterAndSlaveReplication() {
	split := strings.Split(*replicaof, " ")
	address := fmt.Sprintf("%s:%s", split[0], split[1])
	// 给主机发送一个 PING 命令
	master, err := net.Dial("tcp", address)
	if err != nil {
		logger.Error("error: ", err.Error())
	}

	_, err = master.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	//writer := NewWriter(master)
	//var arrays []Value
	//arrays = append(arrays, Value{typ: STRING, bulk: "PING"})
	//err = writer.Write(Value{typ: ARRAY, array: arrays})
	if err != nil {
		logger.Error("error:handshake with master:", err.Error())
	}
}

func configGet(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'config get' command"}
	}
	cmd := args[0].bulk
	if strings.ToUpper(cmd) != "GET" {
		return Value{typ: ERROR, str: "ERR unknown command '" + cmd + "'"}
	}
	key := args[1].bulk
	ConfigsMu.RLock()
	value, ok := Configs[key]
	defer ConfigsMu.RUnlock()
	if !ok {
		return Value{typ: NULL}
	}
	var values []Value
	values = append(values, Value{typ: BULK, bulk: key})
	values = append(values, Value{typ: BULK, bulk: value})
	return Value{typ: ARRAY, array: values}
}
