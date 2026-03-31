package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var directory string = "./"
var port string = "4221"

const workers = 100

const BUFFER_SIZE = 1048576

func main() {

	argsParser()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Logs will appear here")
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println("port binding error")
		os.Exit(1)
	}

	defer listener.Close()

	connections := make(chan net.Conn, 1000)
	for i := 0; i < workers; i++ {
		go deployWorker(context.Background(), connections)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println("Failed to accept Connection")
					continue
				}
			}
			fmt.Println("New Connection established", conn.RemoteAddr().String())
			select {
			case connections <- conn:
			case <-ctx.Done():
				conn.Close()
				return
			}
		}
	}()
	<-ctx.Done()
	fmt.Println()
	fmt.Println(context.Cause(ctx))
	close(connections)
	time.Sleep(time.Second)
	fmt.Println("shutting down complete")
}

func deployWorker(ctx context.Context, connections chan net.Conn) {
	for {
		select {
		case conn, err := <-connections:
			if !err {
				return
			}
			handleConnection(conn)
		case <-ctx.Done():
			return
		}
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		fmt.Println("Connection closed", conn.RemoteAddr().String())
		conn.Close()
	}(conn)

	for {
		readCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		type readResult struct {
			n   int
			err error
		}

		buffer := make([]byte, BUFFER_SIZE)

		readChannel := make(chan readResult, 1)
		go func(conn net.Conn) {
			n, err := conn.Read(buffer)
			readChannel <- readResult{n, err}
		}(conn)

		var n int
		var err error

		select {
		case readResult := <-readChannel:
			n, err = readResult.n, readResult.err
			cancel()
		case <-readCtx.Done():
			cancel()
			fmt.Println("read timepout")
			return
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error", err)
			}
			break
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
				directory = path.Dir(os.Args[i+1])
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
