package main

import (
	"fmt"
	"io"
	"log"
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

		log.Printf("command: \n%s", buffer[:n])

		_, err = conn.Write([]byte("+PONG\r\n"))

		if err != nil {
			fmt.Println("error in write response")
			return
		}
	}
}
