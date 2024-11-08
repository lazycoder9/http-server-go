package main

import (
	"fmt"
	"net"
	"os"
  "strings"
)

func parseRequest(data string) (verb, path string) {
  requestInfo := strings.Split(data, "\r\n")[0]
  parts := strings.Split(requestInfo, " ")
  method := parts[0]
  path = parts[1]

	return method, path
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	_, path := parseRequest(string(buf[:n]))
	switch path {
	case "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Listening on port 4221")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		handleClient(conn)
	}
}
