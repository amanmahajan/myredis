package main

import (
	"flag"
	"log"
	"myredis/config"
	"myredis/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the main ")
	flag.IntVar(&config.Port, "port", 7379, "port for the main")
	flag.Parse()
}
func main() {
	setupFlags()
	log.Println("Starting main")
	err := server.RunAsyncTcpServer()
	if err != nil {
		log.Fatal(err)
	}
}
