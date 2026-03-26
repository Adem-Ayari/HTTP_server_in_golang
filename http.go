package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
)

var directory string = "default_directory"
var port string = "4221"

func main() {

	argsParser()

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println("port binding error")
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Connection acceptance error")
			continue
		}

		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		fmt.Println("Connection closed", conn.RemoteAddr().String())
		conn.Close()
	}(conn)

	for {
		buffer := make([]byte, 1048576)
		n, err := conn.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error", err)
			}
			return
		}

		requestParse := parser(string(buffer[:n]))

		for key, value := range requestParse {
			fmt.Println(key, value)
		}

		response := response(directory, requestParse)
		conn.Write([]byte(response))
		if state, ok := requestParse["Connection"]; ok && state == "close" {
			break
		}
	}
}
func argsParser() {
	if len(os.Args)%2 == 1 {
		for i := 1; i < len(os.Args); i = i + 2 {
			switch os.Args[i] {
			case "--directory":
				directory = path.Dir(os.Args[i+1])
			case "-d":
				directory = os.Args[i+1]
			case "--port":
				port = os.Args[i+1]
			case "-p":
				port = os.Args[i+1]
			default:
				log.Fatal(os.ErrInvalid)
			}
		}
	}
}
