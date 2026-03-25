package main

import (
	"fmt"
	"path"
	"strings"
)

func response(dir string, request map[string]string) string {
	requestLineArgs := strings.Split(request["request line"], " ")
	switch requestLineArgs[0] {
	case "GET":
		return GET(dir, request)
	case "POST":
		return POST(dir, request)
	default:
		return "HTTP/1.1 405 Method Not Found\r\n\r\n"
	}
}

func GET(dir string, request map[string]string) string {

	var response string
	requestLineArgs := strings.Split(request["request line"], " ")
	requestPath := path.Clean(requestLineArgs[1])

	if strings.Contains(requestPath, "/echo/") {
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(requestPath[7:]), requestPath[7:])
		return response
	}

	switch requestPath {
	case "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	case "/user-agent":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request["User-Agent"]), request["User-Agent"])
	default:
		if fileExistance(dir, requestPath) {
			n, body := fileHandlerGET(dir, requestPath)
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", n, body)
		} else {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	}

	return response
}

func POST(dir string, request map[string]string) string {

	var response string
	requestLineArgs := strings.Split(request["request line"], " ")
	requestPath := path.Clean(string(requestLineArgs[1]))

	if pathExistance(dir, requestPath) {
		fileHandlerPOST(dir, requestPath, request["request body"])
		response = "HTTP/1.1 201 Created\r\n\r\n"
	}

	return response
}
