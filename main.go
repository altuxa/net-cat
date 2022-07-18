package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	Conn net.Conn
	Name string
}

type Handler struct {
	clients  map[string]Client
	leaving  chan message
	messages chan message
}

type message struct {
	text    string
	address string
}

func main() {
	arg := os.Args[1:]
	port := "8989"
	if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else if len(arg) == 1 {
		port = arg[0]
		if !IsNumber(port) {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
	}
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on the port :" + port)
	handler := Handler{
		clients:  make(map[string]Client),
		leaving:  make(chan message),
		messages: make(chan message),
	}

	go handler.broadcaster()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handler.handle(conn)
	}
}

func fileRead(filename string) []byte {
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	data, _ := io.ReadAll(file)
	return data
}

func fileWrite(filename string, data string) {
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(data)
}

func (h *Handler) handle(conn net.Conn) {
	logo := fileRead("logo.txt")
	conn.Write(logo)
	conn.Write([]byte("\n"))

	conn.Write([]byte("[ENTER YOUR NAME]: "))

	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)

	if len(clientName) == 0 {
		conn.Close()
		return
	}
	for {
		if len(clientName) != 0 {
			break
		}
	}

	logData := fileRead("log.txt")
	conn.Write(logData)
	conn.Write([]byte("\n"))

	client := Client{
		Conn: conn,
		Name: clientName,
	}

	h.clients[conn.RemoteAddr().String()] = client
	h.messages <- newMessage("\n"+clientName, " has joined our chat...", conn)
	// write 1 time
	conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + clientName + "]" + ":"))
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msgTime := time.Now().Format("2006-01-02 15:04:05")
		msg := "\n" + "[" + msgTime + "]" + "[" + clientName + "]"
		// write to 1 client
		conn.Write([]byte("[" + msgTime + "]" + "[" + clientName + "]" + ":"))
		if len(input.Text()) == 0 {
			continue
		}
		h.messages <- newMessage(msg, ":"+input.Text(), conn)
	}
	// Delete client form map
	delete(h.clients, conn.RemoteAddr().String())

	h.leaving <- newMessage("\n"+clientName, " has left our chat...", conn)

	conn.Close() // ignore errors
}

func newMessage(name, msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    name + msg,
		address: addr,
	}
}

func (h *Handler) broadcaster() {
	for {
		select {
		case msg := <-h.messages:
			fileWrite("log.txt", msg.text)
			for _, client := range h.clients {
				if msg.address == client.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(client.Conn, msg.text) // NOTE: ignoring network errors
				client.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + client.Name + "]" + ":"))
			}

		case msg := <-h.leaving:
			fileWrite("log.txt", msg.text)
			for _, client := range h.clients {
				fmt.Fprintln(client.Conn, msg.text) // NOTE: ignoring network errors
				client.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + client.Name + "]" + ":"))
			}

		}
	}
}

func IsNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
