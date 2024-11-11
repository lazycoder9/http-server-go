package main

import (
	"net"
	"strings"
)

type Middleware func(conn net.Conn, request *Request)

type Middlewares struct {
	middlewares []Middleware
}

func InitMiddlewares() *Middlewares {
	return &Middlewares{[]Middleware{}}
}

func (m *Middlewares) ApplyMiddlewares(conn net.Conn, request *Request) {
	for _, middleware := range m.middlewares {
		middleware(conn, request)
	}
}

func (m *Middlewares) AddMiddleware(mid Middleware) {
	m.middlewares = append(m.middlewares, mid)
}

func ParseHeaders(conn net.Conn, request *Request) {
	raw := request.raw
	requestParts := strings.Split(raw, "\r\n")
	headers := requestParts[1 : len(requestParts)-2]
	parsedHeaders := make(map[string]string)

	for _, header := range headers {
		headerParts := strings.Split(header, ": ")
		headerName := headerParts[0]
		headerValue := headerParts[1]

		parsedHeaders[headerName] = headerValue
	}

	request.headers = parsedHeaders
}
