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

const (
	status200       = "HTTP/1.1 200 OK"
	status201       = "HTTP/1.1 201 Created"
	status404       = "HTTP/1.1 404 Not Found"
	contentTypeText = "Content-Type: text/plain"
	encondingGzip   = "Content-Encoding: gzip"
)

type Handler func(conn net.Conn, request *Request)

func Handle404(conn net.Conn, request *Request) {
	conn.Write([]byte(status404 + "\r\n\r\n"))
}

func HandleHome(conn net.Conn, request *Request) {
	conn.Write([]byte(status200 + "\r\n\r\n"))
}

func HandleUserAgent(conn net.Conn, request *Request) {
	responseBody := request.headers["User-Agent"]

	response := buildResponse(status200, contentTypeText, responseBody)

	conn.Write([]byte(response))
}

func HandleEcho(conn net.Conn, request *Request) {
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
	responseBody := request.params["pathParam"]

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

func HandleGetFiles(conn net.Conn, request *Request) {
	root := os.Args[2]
	fileContent, err := os.ReadFile(root + request.params["pathParam"])

	var response string

	if err != nil {
		response = buildResponse(status404, "", "")
	} else {
		response = buildResponse(status200, "Content-Type: application/octet-stream", string(fileContent))
	}

	conn.Write([]byte(response))
}

func HandlePostFiles(conn net.Conn, request *Request) {
	root := os.Args[2]
	fileName := request.params["pathParam"]

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
