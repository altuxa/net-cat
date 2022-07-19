package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/altuxa/net-cat/internal/helpers"
)

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

func NewHandler() *Handler {
	return &Handler{
		clients:  make(map[string]Client),
		leaving:  make(chan message),
		messages: make(chan message),
	}
}

func (h *Handler) Handle(conn net.Conn) {
	logo := helpers.FileRead("logo.txt")
	fmt.Fprintf(conn, "%s\n[ENTER YOUR NAME]: ", logo)

	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)

	if len(clientName) == 0 {
		fmt.Fprintln(conn, "Try again, name is required")
		conn.Close()
		return
	}

	if len(h.clients) == 10 {
		fmt.Fprintln(conn, "Chat is full, please try again later")
		conn.Close()
		return
	}

	logData := h.LogsReader()
	fmt.Fprintf(conn, "%s", logData)

	client := Client{
		Conn: conn,
		Name: clientName,
	}
	// Add client to the map
	h.clients[conn.RemoteAddr().String()] = client
	h.messages <- newMessage("\n"+clientName, " has joined our chat...", conn)
	conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + clientName + "]" + ":"))
	// scan all msg
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msgTime := time.Now().Format("2006-01-02 15:04:05")
		msg := "\n" + "[" + msgTime + "]" + "[" + clientName + "]" + ":"
		conn.Write([]byte("[" + msgTime + "]" + "[" + clientName + "]" + ":"))
		if len(input.Text()) == 0 {
			continue
		}
		h.messages <- newMessage(msg, input.Text(), conn)
	}

	// Delete client from map
	delete(h.clients, conn.RemoteAddr().String())

	h.leaving <- newMessage("\n"+clientName, " has left our chat...", conn)

	conn.Close()
}

func newMessage(name, msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    name + msg,
		address: addr,
	}
}

func (h *Handler) Broadcaster() {
	for {
		select {
		case msg := <-h.messages:
			helpers.FileWrite("log.txt", msg.text)
			h.LogsWriter(strings.TrimSpace(msg.text))
			for _, client := range h.clients {
				if msg.address == client.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(client.Conn, msg.text)
				client.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + client.Name + "]" + ":"))
			}
		case msg := <-h.leaving:
			helpers.FileWrite("log.txt", msg.text)
			h.LogsWriter(strings.TrimSpace(msg.text))
			for _, client := range h.clients {
				fmt.Fprintln(client.Conn, msg.text)
				client.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + client.Name + "]" + ":"))
			}
		}
	}
}

func (h *Handler) LogsWriter(log string) {
	h.logs = append(h.logs, log+"\n")
}

func (h *Handler) LogsReader() (logs string) {
	for _, s := range h.logs {
		logs = logs + s
	}
	return logs
}
