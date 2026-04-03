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
		if fileExistence(dir, File) {
			n, body := fileHandlerGET(dir, File)
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", acceptType(request, File), n, string(body))
			if acceptEncoding, ok := request["Accept-Encoding"]; ok {
				encodings := strings.Split(acceptEncoding, ", ")
				for _, encoding := range encodings {
					if encoding == "gzip" {
						n, body = encodedGzip(n, body)
						response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", acceptType(request, File), n, string(body))
					}
				}
			}
		} else {
			fmt.Println("hello")
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	case "/user-agent":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request["User-Agent"]), request["User-Agent"])
	default:
		fmt.Println(requestPath)
		if fileExistence(dir, requestPath) {
			n, body := fileHandlerGET(dir, requestPath)
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", acceptType(request, requestPath), n, string(body))
			if acceptEncoding, ok := request["Accept-Encoding"]; ok {
				encodings := strings.Split(acceptEncoding, ", ")
				for _, encoding := range encodings {
					if encoding == "gzip" {
						n, body = encodedGzip(n, body)
						response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", acceptType(request, requestPath), n, string(body))
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
	requestPath := path.Clean(requestLineArgs[1])

	if pathExistence(dir, requestPath) {
		fileHandlerPOST(dir, requestPath, request["request body"])
		response = "HTTP/1.1 201 Created\r\n\r\n"
	}

	return response
}

func acceptType(request map[string]string, requestFile string) string {
	ext := path.Ext(requestFile)

	if ext == "" {
		return "text/plain"
	}

	mimeType := fileExistenceType(ext)

	acceptHeader, ok := request["Accept"]
	if !ok {
		return mimeType
	}

	acceptTypes := parseAcceptHeader(acceptHeader)

	for _, accepted := range acceptTypes {
		if accepted == mimeType || accepted == "*/*" {
			return mimeType
		}
		if prefix, found := strings.CutSuffix(accepted, "/*"); found {
			if strings.HasPrefix(mimeType, prefix+"/") {
				return mimeType
			}
		}
	}

	return "text/plain"
}

func parseAcceptHeader(acceptHeader string) []string {

	acceptHeader = strings.ReplaceAll(acceptHeader, " ", "")

	var result []string

	parts := strings.Split(acceptHeader, ",")

	for _, part := range parts {
		typeAndQ := strings.Split(part, ";")
		mimeType := typeAndQ[0]

		if mimeType != "" {
			result = append(result, mimeType)
		}
	}

	return result

}
