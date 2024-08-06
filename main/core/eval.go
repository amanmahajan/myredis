package core

import (
	"errors"
	"io"
	"log"
)

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
	default:
		return evaluatePing(cmd.Args, c)
	}
}
