package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("port binding error")
		os.Exit(1)
	}

	defer listener.Close()

	buffer := make([]byte, 1024)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Connection acceptance error")
			os.Exit(1)
		}

		defer conn.Close()

		n, err := conn.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error")
			}
			break
		}

		requestParse := parser(string(buffer[:n]))

		for key, value := range requestParse {
			fmt.Println(key, value)
		}

		requestLineArgs := strings.Split(requestParse["request line"], " ")

		if requestLineArgs[1] == "/user-agent" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\nContent-Length: " +
				strconv.Itoa(len(requestParse["User-Agent"])) + "\r\n\r\n" + requestParse["User-Agent"]))
		} else if strings.Contains(requestLineArgs[1], "/echo") {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\nContent-Length: " +
				strconv.Itoa(len(requestLineArgs[1][6:])) + "\r\n\r\n" + requestLineArgs[1][6:]))
		} else if fileExistance(requestLineArgs[1]) {
			contentLength, contentBody := fileHandler(requestLineArgs[1])
			conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\nContent-Length: " +
				strconv.Itoa(contentLength) + "\r\n\r\n" + string(contentBody)))

		} else if requestLineArgs[1] == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\nConnection: close\r\n\r\n"))
		}

		// if strings.Contains(string(buffer[:n]), "GET") {
		// 	args := strings.Split(string(buffer[:n]), " ")
		// 	requestPath := args[1]
		//
		// 	if strings.Compare(requestPath[:5], "/echo") == 0 {
		// 		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\nContent-Length: " +
		// 			strconv.Itoa((len(requestPath[6:]))) + "\r\n\r\n" + requestPath[6:] + "\r\n"))
		// 	} else if strings.Compare(requestPath[:5], "/") == 0 {
		// 		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		// 	} else {
		// 		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		// 	}
		//
		// }
	}
}
