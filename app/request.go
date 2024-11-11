package main

type Request struct {
	method  string
	path    string
	headers map[string]string
	params  map[string]string
	body    string
	raw     string
}

func NewRequest() *Request {
	return &Request{}
}
