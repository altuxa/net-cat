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
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
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

// var (
// 	clients  = make(map[string]Client)
// 	leaving  = make(chan message)
// 	messages = make(chan message)
// )

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
		h.messages <- newMessage(msg, ": "+input.Text(), conn)
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
			for _, conn := range h.clients {
				if msg.address == conn.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn.Conn, msg.text) // NOTE: ignoring network errors
				conn.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + conn.Name + "]" + ":"))
			}

		case msg := <-h.leaving:
			fileWrite("log.txt", msg.text)
			for _, conn := range h.clients {
				fmt.Fprintln(conn.Conn, msg.text) // NOTE: ignoring network errors
				conn.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + conn.Name + "]" + ":"))
			}

		}
	}
}

// package main

// import "net"
// import "fmt"
// import "bufio"
// import "strings" // требуется только ниже для обработки примера

// func main() {

//   fmt.Println("Launching server...")

//   // Устанавливаем прослушивание порта
//   ln, _ := net.Listen("tcp", ":8081")

//   // Открываем порт
//   conn, _ := ln.Accept()

//   // Запускаем цикл
//   for {
//     // Будем прослушивать все сообщения разделенные \n
//     message, _ := bufio.NewReader(conn).ReadString('\n')
//     // Распечатываем полученое сообщение
//     fmt.Print("Message Received:", string(message))
//     // Процесс выборки для полученной строки
//     newmessage := strings.ToUpper(message)
//     // Отправить новую строку обратно клиенту
//     conn.Write([]byte(newmessage + "\n"))
//   }
// }

// func Socket() {
// 	// Подключаемся к сокету
// 	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
// 	for {
// 		// Чтение входных данных от stdin
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("Text to send: ")
// 		text, _ := reader.ReadString('\n')
// 		// Отправляем в socket
// 		fmt.Fprintf(conn, text+"\n")
// 		// Прослушиваем ответ
// 		message, _ := bufio.NewReader(conn).ReadString('\n')
// 		fmt.Print("Message from server: " + message)
// 	}
// }

/////////////////////////////////////////////////////////
// package server

// import (
// 	"bufio"
// 	"log"
// 	"net"
// 	"strings"
// )

// type Server struct {
// 	Addr string
// }

// func (srv Server) ListenAndServe() error {
// 	addr := srv.Addr
// 	if addr == "" {
// 		addr = ":8080"
// 	}
// 	log.Printf("starting server on %v\n", addr)
// 	listener, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer listener.Close()
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("error accepting connection %v", err)
// 			continue
// 		}
// 		log.Printf("accepted connection from %v", conn.RemoteAddr())
// 		handle(conn)
// 	}
// }

// func handle(conn net.Conn) error {
// 	defer func() {
// 		log.Printf("closing connection from %v", conn.RemoteAddr())
// 		conn.Close()
// 	}()
// 	r := bufio.NewReader(conn)
// 	w := bufio.NewWriter(conn)
// 	scanr := bufio.NewScanner(r)
// 	for {
// 		scanned := scanr.Scan()
// 		if !scanned {
// 			if err := scanr.Err(); err != nil {
// 				log.Printf("%v(%v)", err, conn.RemoteAddr())
// 				return err
// 			}
// 			break
// 		}
// 		w.WriteString(strings.ToUpper(scanr.Text()) + "\n")
// 		w.Flush()
// 	}
// 	return nil
// }
