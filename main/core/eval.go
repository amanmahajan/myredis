package core

import (
	"bytes"
	"errors"
	"fmt"
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
		case "BGREWRITEAOF":
			buffer.Write(evaluateBGREWRITEAOF(cmd.Args))
		case "INCR":
			buffer.Write(evaluateIncr(cmd.Args))
		case "INFO":
			buffer.Write(evaluateInfo(cmd.Args))
		case "CLIENT":
			buffer.Write(evaluateClient())
		case "LATENCY":
			buffer.Write(evaluateLatency())
		case "LRU":
			buffer.Write(evaluateLRU())
		default:
			buffer.Write(evaluatePing(cmd.Args))
		}
	}
	c.Write(buffer.Bytes())
}

func evaluateBGREWRITEAOF(args []string) []byte {
	DumpData()
	return RespOk

}

func evaluateLRU() []byte {
	evictAllkeysLRU()
	return RespOne

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
		return Encode(errors.New("Not enough arguments for the set command"), false)
	}

	var key, val string
	var duration int64 = -1

	key = args[0]
	val = args[1]

	objType, objEnc := deduceTypeEncoding(val)

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

	Put(key, NewObject(val, duration, objType, objEnc))
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

	if hasExpired(obj) {
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

	exp, ok := getExpiry(obj)

	if !ok {
		return NilResp
	}

	if exp <= uint64(time.Now().UnixMilli()) {
		return NoKeyExist
	}

	remainingDurationMs := exp - uint64(time.Now().UnixMilli())

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

func evaluateIncr(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("Wrong arguments for incr"), false)
	}

	key := args[0]
	val := Get(key)
	if val == nil {
		val = NewObject("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		Put(key, val)
	}
	if err := assertType(val.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode(err, false)
	}

	if err := assertEncoding(val.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode(err, false)
	}

	i, _ := strconv.ParseInt(val.Value.(string), 10, 64)
	i++
	val.Value = strconv.FormatInt(i, 10)

	return Encode(i, false)

}

func evaluateInfo(args []string) []byte {
	info := make([]byte, 0)
	buffer := bytes.NewBuffer(info)
	buffer.WriteString("# Keyspace\r\n")
	for i := range KeySpaceStats {
		buffer.WriteString(fmt.Sprintf("db%d:keys=%d,expires=0,avg_ttl=0\r\n", i, KeySpaceStats[i]["keys"]))
	}
	return Encode(buffer.String(), false)

}

func evaluateLatency() []byte {
	return Encode([]string{}, false)
}

func evaluateClient() []byte {
	return RespOk
}
