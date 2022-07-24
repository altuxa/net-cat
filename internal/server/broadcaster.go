package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/altuxa/net-cat/internal/helpers"
)

func (h *Handler) Broadcaster() {
	for {
		select {
		case msg := <-h.messages:
			err := helpers.FileWrite("log.txt", msg.text)
			if err != nil {
				log.Println(err)
			}
			h.LogsWriter(strings.TrimSpace(msg.text))
			msgTime := time.Now().Format("2006-01-02 15:04:05")
			h.mut.Lock()
			for _, client := range h.clients {
				if msg.address == client.Conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(client.Conn, msg.text)
				client.Conn.Write([]byte(fmt.Sprintf("[%s][%s]:", msgTime, client.Name)))
			}
			h.mut.Unlock()
		case msg := <-h.leaving:
			err := helpers.FileWrite("log.txt", msg.text)
			if err != nil {
				log.Println(err)
			}
			h.LogsWriter(strings.TrimSpace(msg.text))
			msgTime := time.Now().Format("2006-01-02 15:04:05")
			h.mut.Lock()
			for _, client := range h.clients {
				fmt.Fprintln(client.Conn, msg.text)
				client.Conn.Write([]byte(fmt.Sprintf("[%s][%s]:", msgTime, client.Name)))
			}
			h.mut.Unlock()
		}
	}
}
