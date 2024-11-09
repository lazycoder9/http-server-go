package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	status200       = "HTTP/1.1 200 OK"
	status404       = "HTTP/1.1 404 Not Found"
	contentTypeText = "Content-Type: text/plain"
)

func handleEcho(conn net.Conn, path string) {
	responseBody := path[6:]

	response := buildResponse(status200, contentTypeText, responseBody)

	conn.Write([]byte(response))
}

func handleUserAgent(conn net.Conn, headers map[string]string) {
	responseBody := headers["User-Agent"]

	response := buildResponse(status200, contentTypeText, responseBody)

	conn.Write([]byte(response))
}

func parseRequest(data string) (method, path string, headers map[string]string) {
	requestParts := strings.Split(data, "\r\n")
	requestStatusLine := requestParts[0]
	parts := strings.Split(requestStatusLine, " ")
	method = parts[0]
	path = parts[1]

	headers = parseHeaders(requestParts[1 : len(requestParts)-2])

	return method, path, headers
}

func parseHeaders(headers []string) map[string]string {
	parsedHeaders := make(map[string]string)

	for _, header := range headers {
		headerParts := strings.Split(header, ": ")
		headerName := headerParts[0]
		headerValue := headerParts[1]

		parsedHeaders[headerName] = headerValue
	}

	return parsedHeaders
}

func buildResponse(status, contentTypeHeader, body string) string {
	contentLengthHeader := fmt.Sprintf("Content-Length: %d", len(body))

	responseHeaders := strings.Join([]string{
		status,
		contentTypeHeader,
		contentLengthHeader,
	}, "\r\n")

	return responseHeaders + "\r\n\r\n" + body
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	_, path, headers := parseRequest(string(buf[:n]))

	switch {
	case path == "/":
		conn.Write([]byte(status200 + "\r\n\r\n"))
	case path == "/user-agent":
		handleUserAgent(conn, headers)
	case strings.HasPrefix(path, "/echo"):
		handleEcho(conn, path)
	default:
		conn.Write([]byte(status404 + "\r\n\r\n"))
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

		go handleClient(conn)
	}
}
