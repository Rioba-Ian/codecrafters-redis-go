package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Type byte

type RESP struct {
	Type Type
	Raw  []byte
	Data []byte
}

/*
 * Resp data type | Category | first byte
 * Simple strings | Simple |  +
 * Simple Errors | Simple |  -
 * Simple Integers | Simple | :
 * Bulk Strings | Aggregate | $
 * Array strings |Aggregate | *
 */

func (r *RESP) parse() (string, []byte) {
	splitString := strings.Split(string(r.Raw), "\r\n")
	packets, _ := json.Marshal(strings.Join(splitString, ""))
	// fmt.Println("\nraw data::", string(packets))

	for i := 0; i < len(packets); i++ {
		ch := packets[i]

		if ch == '$' {
			size := readInt(packets[i+1])
			takeSlice := packets[i+2 : i+2+size]
			r.Data = append(r.Data, takeSlice...)
			r.Data = append(r.Data, '%')
		}

	}

	cmd := strings.Split(string(r.Data), "%")

	if len(cmd) > 1 {
		byteRes := addPrefix(cmd[1])
		return cmd[0], byteRes
	}

	return cmd[0], nil
}

func readInt(b byte) int {
	x, _ := strconv.Atoi(string(b))
	return x
}

func addPrefix(input string) []byte {
	finalBytes := []byte("$")
	inputBytes, _ := json.Marshal(input)
	lenBytes, _ := json.Marshal(len(input))
	firstBytes := []byte("\r\n")

	finalBytes = append(finalBytes, lenBytes...)
	finalBytes = append(finalBytes, firstBytes...)

	return append(finalBytes, inputBytes...)
}
