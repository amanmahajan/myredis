package core

import (
	"errors"
	"io"
	"log"
	"myredis/server"
	"strconv"
	"time"
)

var NilResp = []byte("$-1\r\n")
var NoKeyExist = []byte("$-2\r\n")

func evaluatePing(args []string, conn io.ReadWriter) error {
	byteArr := make([]byte, 0)
	if len(args) >= 2 {
		return errors.New("Not enough arguments")
	}
	if len(args) == 0 { // PING Scenario
		byteArr = Encode("PING", true)

	} else if len(args) == 1 { // Ping "Hello" scenario
		byteArr = Encode(args[0], true)
	}

	_, err := conn.Write(byteArr)
	if err != nil {
		return err
	}
	return nil

}
func EvalAndRespond(cmd *RedisCommand, c io.ReadWriter) error {
	log.Println("comamnd:", cmd.Command)
	switch cmd.Command {
	case "PING":
		return evaluatePing(cmd.Args, c)
	case "SET":
		return evaluateSet(cmd.Args, c)
	case "GET":
		return evaluateGet(cmd.Args, c)
	case "TTL":
		return evaluateTTL(cmd.Args, c)
	default:
		return evaluatePing(cmd.Args, c)
	}
}

func evaluateSet(args []string, conn io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("Not enough arguments")
	}

	var key, val string
	var duration int64 = -1

	for k := 2; k < len(args); k++ {

		switch args[k] {
		case "EX", "ex":
			k++
			if len(args) <= k {
				return errors.New("Not enough arguments for expiration")
			}
			timeSec, err := strconv.ParseInt(args[k], 10, 64)

			if err != nil {
				return errors.New("Invalid expiration")
			}
			duration = timeSec * 1000
		default:
			return errors.New("Invalid command")
		}
	}

	server.Put(key, server.NewObject(val, duration))
	conn.Write([]byte("+OK\r\n")) // Sending OK output
	return nil
}

func evaluateGet(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("Wrong arguments for get")
	}

	key := args[0]

	obj := server.Get(key)
	if obj == nil {
		conn.Write(NilResp)
		return nil
	}

	if obj.Expiry != -1 && obj.Expiry < time.Now().UnixMilli() {
		conn.Write(NilResp)
		return nil
	}
	conn.Write(Encode(obj.Value, false))
	return nil
}

// Get the ttl value from the key
func evaluateTTL(args []string, conn io.ReadWriter) error {

	if len(args) != 1 {
		return errors.New("Wrong arguments for evaluating ttl")
	}

	key := args[0]

	obj := server.Get(key)
	if obj == nil {
		conn.Write(NoKeyExist)
		return nil
	}
	if obj.Expiry == -1 {
		conn.Write(NilResp)
		return nil
	}

	remainingDurationMs := -obj.Expiry - time.Now().UnixMilli()
	if remainingDurationMs < 0 {
		conn.Write(NoKeyExist)
		return nil

	}
	c.Write(Encode(int64(remainingDurationMs/1000), false))
	return nil

}
