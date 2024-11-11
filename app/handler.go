package main

import "net"

type Handler func(conn net.Conn, request *Request)
