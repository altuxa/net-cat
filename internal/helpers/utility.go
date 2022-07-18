package helpers

import (
	"errors"
	"io"
	"os"
)

func IsNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func CheckPort(arg []string) (string, error) {
	port := "8989"
	if len(arg) > 1 {
		return "", errors.New("[USAGE]: ./TCPChat $port")
	} else if len(arg) == 1 {
		if !IsNumber(arg[0]) {
			return "", errors.New("[USAGE]: ./TCPChat $port")
		}
		port = arg[0]
	}
	return port, nil
}

func FileRead(filename string) []byte {
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	data, _ := io.ReadAll(file)
	return data
}

func FileWrite(filename string, data string) {
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(data)
}
