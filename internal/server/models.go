package server

import (
	"net"
	"sync"
)

type Handler struct {
	clients  map[string]Client
	leaving  chan message
	messages chan message
	logs     []string
	mut      sync.Mutex
}

type Client struct {
	Conn net.Conn
	Name string
}

type message struct {
	text    string
	address string
}
