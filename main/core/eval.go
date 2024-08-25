package core

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

var NilResp = []byte(":-1\r\n")
var NoKeyExist = []byte(":-2\r\n")
var RespOk = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")

func EvalAndRespond(cmds RedisCommands, c io.ReadWriter) {

	byteArr := make([]byte, 0)
	buffer := bytes.NewBuffer(byteArr)
	for _, cmd := range cmds {

		log.Println("comamnd:", cmd.Command)
		switch cmd.Command {
		case "PING":
			buffer.Write(evaluatePing(cmd.Args))
		case "SET":

			buffer.Write(evaluateSet(cmd.Args))
		case "GET":

			buffer.Write(evaluateGet(cmd.Args))
		case "TTL":
			buffer.Write(evaluateTTL(cmd.Args))
		case "DEL":

			buffer.Write(evaluateDelete(cmd.Args))
		case "EXPIRE":

			buffer.Write(evaluateExpire(cmd.Args))

		default:
			buffer.Write(evaluatePing(cmd.Args))
		}
	}
	c.Write(buffer.Bytes())
}

func evaluatePing(args []string) []byte {
	byteArr := make([]byte, 0)
	if len(args) >= 2 {
		return Encode(errors.New("Not enough arguments"), false)
	}
	if len(args) == 0 { // PING Scenario
		byteArr = Encode("PONG", true)

	} else if len(args) == 1 { // Ping "Hello" scenario
		byteArr = Encode(args[0], true)
	}

	return byteArr

}

func evaluateSet(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("Not enough arguments"), false)
	}

	var key, val string
	var duration int64 = -1

	key = args[0]
	val = args[1]

	for k := 2; k < len(args); k++ {

		switch args[k] {
		case "EX", "ex":
			k++
			if len(args) <= k {
				return Encode(errors.New("Not enough arguments for expiration"), false)
			}
			timeSec, err := strconv.ParseInt(args[k], 10, 64)

			if err != nil {
				return Encode(errors.New("Invalid expiration"), false)
			}
			duration = timeSec * 1000
		default:
			return Encode(errors.New("Invalid command"), false)
		}
	}

	Put(key, NewObject(val, duration))
	return RespOk
}

func evaluateGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("Wrong arguments for get"), false)
	}

	key := args[0]

	obj := Get(key)
	if obj == nil {
		return NilResp
	}

	if obj.Expiry != -1 && obj.Expiry < time.Now().UnixMilli() {
		return NilResp
	}
	return Encode(obj.Value, false)
}

// Get the ttl value from the key
func evaluateTTL(args []string) []byte {

	if len(args) != 1 {
		return Encode(errors.New("Wrong arguments for evaluating ttl"), false)
	}

	key := args[0]

	obj := Get(key)
	if obj == nil {
		return NoKeyExist
	}
	if obj.Expiry == -1 {
		return NoKeyExist
	}

	remainingDurationMs := -obj.Expiry - time.Now().UnixMilli()
	if remainingDurationMs < 0 {
		return NoKeyExist

	}
	return Encode(int64(remainingDurationMs/1000), false)
}

func evaluateDelete(args []string) []byte {
	totalDel := 0
	for _, str := range args {
		if Delete(str) {
			totalDel++
		}

	}
	return Encode(totalDel, false)
}

func evaluateExpire(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("Wrong arguments for evaluating expire"), false)
	}

	keyStr := args[0]
	expireTime, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(err, false)
	}

	key := Get(keyStr)
	if key == nil {
		// Sending 0 as timeout not set and key does not exist
		return RespZero
	}

	key.Value = time.Now().UnixMilli() + expireTime*1000
	return RespOne
}
