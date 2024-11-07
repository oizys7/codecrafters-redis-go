package main

import "github.com/codecrafters-io/redis-starter-go/logging"

var logger = logging.New(logging.LevelDebug)

func main() {
	initConfigs()
	loadRdbFileIntoKVMemoryStore()

	server := &Server{}
	defer server.Close()
	server.Start()
}
