package main

import (
	"fmt"
	"net"
)

type Router struct {
	routes map[string]map[string]func(net.Conn, *Request)
}

func NewRouter() *Router {
	routes := make(map[string]map[string]func(net.Conn, *Request))
	return &Router{routes}
}

func (router *Router) addRoute(method, path string, handler func(net.Conn, *Request)) {
	if _, exists := router.routes[method]; !exists {
		router.routes[method] = make(map[string]func(net.Conn, *Request))
	}

	router.routes[method][path] = handler
}

func (router *Router) listRoutes() {
	fmt.Println("List of routes")
	routes := router.routes
	for method, routes := range routes {
		for path, _ := range routes {
			fmt.Println(method, path)
		}
	}
}

func (router *Router) route() {
	fmt.Println("Hello from Router!")
}
