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

/*
Simple read method to read 512 bytes from the connection at a time

TODO : Try to implement read more than 512 bytes later
*/
func readCommand(c net.Conn) (*core.RedisCommand, error) {

	buf := make([]byte, 512)
	// Read 512 bytes for now. Implement reading more later
	_, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	val, err := core.DecodeArrayString(buf)
	if err != nil {
		return nil, err
	}
	return &core.RedisCommand{
		Command: strings.ToUpper(val[0]),
		Args:    val[1:],
	}, nil

}

/*
Reply to the client
*/
func respond(cmd *core.RedisCommand, c net.Conn) error {

	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
	return nil
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
			cmd, err := readCommand(c)
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
			log.Println("Command received:", cmd)
			if err := respond(cmd, c); err != nil {
				log.Println("Error responding to client:", err)
			}
		}
	}

}

func respondError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
