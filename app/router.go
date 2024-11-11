package main

import "fmt"

type Router struct {
	routes map[string]func(foo string)
}

func NewRouter() *Router {
	routes := make(map[string]func(foo string))
	return &Router{routes}
}

func (router *Router) route() {
	fmt.Println("Hello from Router!")
}
