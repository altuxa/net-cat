package server

import "net"

type Handler struct {
	clients  map[string]Client
	leaving  chan message
	messages chan message
	logs     []string
}

type Client struct {
	Conn net.Conn
	Name string
}

type message struct {
	text    string
	address string
}
