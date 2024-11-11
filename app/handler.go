package main

import "net"

type Handler struct {
	conn    net.Conn
	request Request
}
