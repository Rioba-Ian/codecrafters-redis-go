package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const PONG_RES = "+PONG\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// create byte slice
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by peer")
			} else {
				fmt.Println("Error reading from connection of peer")
			}
			return
		}

		resp := RESP{}

		resp.Raw = buffer[:n]
		cmd, dataBytes := resp.parse()

		switch cmd {
		case "ECHO", "echo":
		case "PING", "ping":
			dataBytes = []byte(PONG_RES)
		default:
		}

		_, err = conn.Write(dataBytes)

		if err != nil {
			fmt.Println("error in write response")
			return
		}
	}
}
