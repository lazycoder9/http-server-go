package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
)


func handleEcho(conn net.Conn, request *Request) {
	supportedEndodings := []string{"gzip"}
	encodings, exists := request.headers["Accept-Encoding"]

	var acceptedEncoding string

	for _, e := range strings.Split(encodings, ", ") {
		if slices.Contains(supportedEndodings, e) {
			acceptedEncoding = e
			break
		}
	}

	responseHeaders := []string{contentTypeText}
	responseBody := request.path[6:]

	if exists && acceptedEncoding == "gzip" {
		responseHeaders = append(responseHeaders, encondingGzip)

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write([]byte(responseBody))

		if err != nil {
			fmt.Println("Error writing to gzip: ", err.Error())
		}

		zw.Close()
		responseBody = buf.String()
	}

	response := buildResponse(status200, strings.Join(responseHeaders, "\r\n"), responseBody)

	conn.Write([]byte(response))
}

func handleGetFiles(conn net.Conn, request *Request) {
	root := os.Args[2]
	fileContent, err := os.ReadFile(root + strings.TrimPrefix(request.path, "/files/"))

	var response string

	if err != nil {
		response = buildResponse(status404, "", "")
	} else {
		response = buildResponse(status200, "Content-Type: application/octet-stream", string(fileContent))
	}

	conn.Write([]byte(response))
}

func handlePostFiles(conn net.Conn, request *Request) {
	root := os.Args[2]
	fileName := strings.TrimPrefix(request.path, "/files/")

	var response string

	file, err := os.Create(root + fileName)

	if err != nil {
		response = buildResponse(status404, "", "")
	}

	_, errWrite := file.Write([]byte(request.body))

	if errWrite != nil {
		response = buildResponse(status404, "", "")
	}

	response = buildResponse(status201, "", "")
	conn.Write([]byte(response))
}

func handleUserAgent(conn net.Conn, request *Request) {
	responseBody := request.headers["User-Agent"]

	response := buildResponse(status200, contentTypeText, responseBody)

	conn.Write([]byte(response))
}

func parseRequest(data string) (method, path string, headers map[string]string, body string) {
	requestParts := strings.Split(data, "\r\n")
	requestStatusLine := requestParts[0]
	parts := strings.Split(requestStatusLine, " ")
	method = parts[0]
	path = parts[1]
	body = requestParts[len(requestParts)-1]

	headers = parseHeaders(requestParts[1 : len(requestParts)-2])

	return method, path, headers, body
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

	method, path, headers, body := parseRequest(string(buf[:n]))
	request := NewRequest(method, path, headers, body)

  handler := router.route(request)
  handler(conn, request)

	// switch {
	// case path == "/":
	// case path == "/user-agent":
	// 	handleUserAgent(conn, request)
	// case strings.HasPrefix(path, "/echo"):
	// 	handleEcho(conn, request)
	// case strings.HasPrefix(path, "/files") && method == "GET":
	// 	handleGetFiles(conn, request)
	// case strings.HasPrefix(path, "/files") && method == "POST":
	// 	handlePostFiles(conn, request)
	// default:
	// 	conn.Write([]byte(status404 + "\r\n\r\n"))
	// }
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

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleRequest(conn, router)
	}
}
