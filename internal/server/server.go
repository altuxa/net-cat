package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/altuxa/net-cat/internal/helpers"
)

func NewHandler() *Handler {
	return &Handler{
		clients:  make(map[string]Client),
		leaving:  make(chan message),
		messages: make(chan message),
	}
}

func (h *Handler) Handle(conn net.Conn) {
	logo, err := helpers.FileRead("logo.txt")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(conn, "%s, oops something went wrong and the logo didn't load\n", err)
	}

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
	h.mut.Lock()
	h.clients[conn.RemoteAddr().String()] = client
	h.mut.Unlock()
	h.messages <- newMessage("\n"+clientName, " has joined our chat...", conn)
	fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), clientName)

	// scan all msg
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msgTime := time.Now().Format("2006-01-02 15:04:05")
		msg := fmt.Sprintf("\n[%s][%s]:", msgTime, clientName)
		fmt.Fprintf(conn, "[%s][%s]:", msgTime, clientName)
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
