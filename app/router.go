package main

import (
	"fmt"
)

type Router struct {
	routes map[string]map[string]Handler
}

func NewRouter() *Router {
	routes := make(map[string]map[string]Handler)
	return &Router{routes}
}

func (router *Router) addRoute(method, path string, handler Handler) {
	if _, exists := router.routes[method]; !exists {
		router.routes[method] = make(map[string]Handler)
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

func (router *Router) route(request *Request) Handler {
  routes := router.routes[request.method]

  handler := fetchHandler(routes, request)

  return handler
}

func fetchHandler(routes map[string]Handler, request *Request) Handler {
  handler, exists := routes[request.path]

  if !exists {
    fmt.Printf("Handler for %s %s does not exist\n", request.method, request.path)
    return Handle404
  }

  return handler
}
