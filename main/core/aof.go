package core

import (
	"fmt"
	"log"
	"myredis/config"
	"os"
	"strings"
)

func dumpKey(fb *os.File, key string, obj *Obj) {

	command := fmt.Sprintf("SET %s %s", key, obj.Value)
	tokens := strings.Split(command, " ")

	_, err := fb.Write(Encode(tokens, false))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DumpData() {
	fl, err := os.OpenFile(config.AOFFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("rewritting AOF file at", config.AOFFile)
	for k, v := range store {
		dumpKey(fl, k, v)
	}
	log.Println("rewritting AOF file at", config.AOFFile)
}
