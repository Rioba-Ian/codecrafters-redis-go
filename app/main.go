package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

/*
 * TODO:- handle rdb config get
 *
 */

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	rdbFileConfig()
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

	parser := NewRediParser(conn)

	for {
		value, err := parser.Parse()

		if err == io.EOF {
			fmt.Println("Connection closed by peer")
			break
		}

		if err != nil {
			fmt.Println("Error parsing command")
			return
		}

		cmd := value.convertToCommand()
		res := handleCommand(cmd)

		_, err = conn.Write([]byte(res))

		if err != nil {
			fmt.Println("error in write response")
			return
		}
	}
}
