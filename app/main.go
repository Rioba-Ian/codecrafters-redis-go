package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

type Command struct {
	echoOutput string
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

		cmd, restStr := parser(string(buffer[:n]))
		var byteRes []byte
		switch cmd {
		case "ECHO", "echo":
			jsonBytes, _ := json.Marshal(strings.Join(restStr, ""))
			byteRes = jsonBytes
		case "PING", "ping":
			byteRes = []byte("+PONG\r\n")
		default:
		}

		_, err = conn.Write(byteRes)

		if err != nil {
			fmt.Println("error in write response")
			return
		}
	}
}

/*
 * Resp data type | Category | first byte
 * Simple strings | Simple |  +
 * Simple Errors | Simple |  -
 * Simple Integers | Simple | :
 * Bulk Strings | Aggregate | $
 * Array strings |Aggregate | *
 */

func parser(input string) (string, []string) {

	words := strings.Fields(input)
	var receivedStr []string

	for _, ch := range words {
		if strings.HasPrefix(ch, "*") {
		} else if strings.HasPrefix(ch, "$") {
		} else {
			receivedStr = append(receivedStr, ch)
		}
	}

	return receivedStr[0], receivedStr[1:]
}
