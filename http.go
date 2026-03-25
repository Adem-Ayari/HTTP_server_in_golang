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
			os.Exit(1)
		}

		defer conn.Close()

		buffer := make([]byte, 1048576)
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

		response := response(directory, requestParse)
		conn.Write([]byte(response))
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
