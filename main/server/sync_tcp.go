package server

import (
	"fmt"
	"io"
	"log"
	"myredis/config"
	"myredis/core"
	"net"
	"strconv"
	"strings"
)

func toArrayString(intfs []interface{}) ([]string, error) {

	values := make([]string, len(intfs))
	for i := range intfs {
		values[i] = intfs[i].(string)
	}
	return values, nil

}

/*
Simple read method to read 512 bytes from the connection at a time

TODO : Try to implement read more than 512 bytes later
*/
func readCommands(c io.ReadWriter) (core.RedisCommands, error) {

	buf := make([]byte, 512)
	n, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	commands := make([]*core.RedisCommand, 0)
	for _, v := range values {
		tokens, err := toArrayString(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		commands = append(commands, &core.RedisCommand{
			Command: strings.ToUpper(tokens[0]),
			Args:    tokens[1:],
		})
	}
	return commands, nil

}

/*
Reply to the client
*/
func respond(cmd core.RedisCommands, c io.ReadWriter) {

	core.EvalAndRespond(cmd, c)
}

// RunSyncTcpServer Inefficient TCP server that accepts only 1 connection at a time
func RunSyncTcpServer() {
	log.Println("Starting sync Tcp Server on ", config.Port, config.Host)

	currentClients := 0

	// starting listening to the port
	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		// Blocking call
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		currentClients++
		log.Println("Client go connected to ", c.RemoteAddr(), "total connections ", currentClients)

		// Client connected. Start reading the data
		for {
			cmds, err := readCommands(c)
			if err != nil {
				c.Close()
				currentClients--
				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients ", currentClients)
				// If clients wants to disconnect
				if err == io.EOF {
					break
				}
				log.Println("err", err)

			}
			respond(cmds, c)

		}
	}

}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
