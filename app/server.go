package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func parseRequest(data string) (method, path string, headers, params map[string]string, body string) {
	requestParts := strings.Split(data, "\r\n")
	requestStatusLine := requestParts[0]
	parts := strings.Split(requestStatusLine, " ")
	method = parts[0]
	path = parts[1]
	body = requestParts[len(requestParts)-1]

	headers = parseHeaders(requestParts[1 : len(requestParts)-2])
	pathParts := strings.Split(path, "/")

	parsedPath := "/" + pathParts[1]
	var pathParam string

	if len(pathParts) > 2 {
		pathParam = pathParts[2]
	}

	params = map[string]string{
		"pathParam": pathParam,
	}

	return method, parsedPath, headers, params, body
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

func handleRequest(conn net.Conn, router *Router) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	method, path, headers, params, body := parseRequest(string(buf[:n]))
	request := NewRequest(method, path, headers, params, body)

	handler := router.route(request)
	handler(conn, request)
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Listening on port 4221")
	router := NewRouter()
	router.addRoute("GET", "/", HandleHome)
	router.addRoute("GET", "/user-agent", HandleUserAgent)
	router.addRoute("GET", "/echo/{value}", HandleEcho)
	router.addRoute("GET", "/files/{filename}", HandleGetFiles)
	router.addRoute("POST", "/files/{filename}", HandlePostFiles)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleRequest(conn, router)
	}
}
