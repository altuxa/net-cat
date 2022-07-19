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

func FileRead(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func FileWrite(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(data)
	return nil
}
