package main

import (
	"fmt"
	"path"
	"strings"
)

var SupportedCompression []string = []string{"gzip"}

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
	Content_Length := 0
	var responseBody []byte

	if strings.Contains(requestPath, "/echo/") {
		Content_Length = len(requestPath[6:])
		responseBody = []byte(requestPath[6:])

		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", Content_Length, string(responseBody))

		if acceptEncoding, ok := request["Accept-Encoding"]; ok {
			encodings := strings.Split(acceptEncoding, ", ")
			for _, encoding := range encodings {
				if encoding == "gzip" {
					Content_Length, responseBody = encodedGzip(Content_Length, responseBody)
					response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", Content_Length, string(responseBody))
				}
			}
		}

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
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", n, string(body))
			if acceptEncoding, ok := request["Accept-Encoding"]; ok {
				encodings := strings.Split(acceptEncoding, ", ")
				for _, encoding := range encodings {
					if encoding == "gzip" {
						n, body = encodedGzip(n, body)
						response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", n, string(body))
					}
				}
			}
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
