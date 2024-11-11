package main

type Request struct {
	method  string
	path    string
	headers map[string]string
	body    string
}

func NewRequest(method, path string, headers map[string]string, body string) *Request {
	return &Request{method, path, headers, body}
}
