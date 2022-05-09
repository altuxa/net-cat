package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port")
)

func main() {
	flag.Parse()
	fmt.Println(*listen)
	if *listen {
		startServer()
		return
	}
	fmt.Println(flag.Arg(0))
	fmt.Println(flag.Arg(1))
	fmt.Println(*listen)
	if len(flag.Args()) < 2 {
		fmt.Println("Hostname and port required")
		return
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	startClient(fmt.Sprintf("%s:%s", serverHost, serverPort))
}

func startServer() {
	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			// fmt.Println(conn)
			var text string
			// fmt.Scan(&text)
			go processClient(text, conn)
		}
	}
}

func processClient(s string, conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	// message, _ := bufio.NewReader(conn).ReadString('\n')
	// fmt.Print(s + message)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
}

func startClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Can't connect to server: %s\n", err)
		return
	}
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("Connection error: %s\n", err)
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
