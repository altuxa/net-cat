package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handle(conn)
	}
}

type Client struct {
	Conn net.Conn
	Name string
}

var (
	clients  = make(map[string]Client)
	leaving  = make(chan message)
	messages = make(chan message)
)

type message struct {
	text    string
	address string
}

func handle(conn net.Conn) {
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	// in := bufio.NewScanner(conn)
	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)
	fmt.Println(clientName)
	if len(clientName) == 0 {
		conn.Close()
		return
	}
	for {
		if len(clientName) != 0 {
			break
		}
	}

	client := Client{
		Conn: conn,
		Name: clientName,
	}

	clients[conn.RemoteAddr().String()] = client
	messages <- newMessage("\n"+clientName, " has joined our chat...", conn)
	// write 1 time
	conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + clientName + "]" + ":"))
	input := bufio.NewScanner(conn)
	// write to all users
	msg := "\n" + "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + clientName + "]"
	conn.Read([]byte(""))
	for input.Scan() {
		// write to 1 client
		conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + clientName + "]" + ":"))
		if len(input.Text()) == 0 {
			continue
		}
		messages <- newMessage(msg, ": "+input.Text(), conn)
	}
	// Delete client form map
	delete(clients, conn.RemoteAddr().String())

	leaving <- newMessage("\n"+clientName, " has left our chat...", conn)

	conn.Close() // ignore errors
}

func newMessage(name, msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    name + msg,
		address: addr,
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				if msg.address == conn.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn.Conn, msg.text) // NOTE: ignoring network errors
				conn.Conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + conn.Name + "]" + ":"))
			}

		case msg := <-leaving:
			for _, conn := range clients {
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
