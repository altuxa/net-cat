package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/altuxa/net-cat/internal/helpers"
	"github.com/altuxa/net-cat/internal/server"
)

func main() {
	arg := os.Args[1:]
	port, err := helpers.CheckPort(arg)
	if err != nil {
		fmt.Println(err)
		return
	}
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on the port :%s\n", port)
	handler := server.NewHandler()
	go handler.Broadcaster()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler.Handle(conn)
	}
}
