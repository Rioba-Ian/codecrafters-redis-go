package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type RESP struct {
	Type   string
	String string
	Number int
	Raw    []byte
	Data   []byte
	Array  []RESP
}

type RedisParser struct {
	reader *bufio.Reader
}

func NewRediParser(reader io.Reader) *RedisParser {
	return &RedisParser{
		reader: bufio.NewReader(reader),
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

func (p *RedisParser) Parse() (*RESP, error) {
	typeByte, err := p.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typeByte {
	case '+':
		return p.parseSimpleString()
	case '-':
		fmt.Println("here is where it should be running")
		return p.parseErrorString()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArrayString()
	default:
		return nil, fmt.Errorf("unknown start of byte: %c", typeByte)
	}
}

func (p *RedisParser) parseSimpleString() (*RESP, error) {
	line, err := p.readLine()

	if err != nil {
		return nil, err
	}

	return &RESP{Type: "string", String: line}, nil
}

func (p *RedisParser) parseErrorString() (*RESP, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	fmt.Printf("The error string:: %s\n", line)

	return &RESP{Type: "error", String: line}, nil
}

func (p *RedisParser) parseArrayString() (*RESP, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(line)

	if err != nil {
		return nil, fmt.Errorf("invalid integer: %s", line)
	}

	if length == -1 {
		return &RESP{Type: "null"}, nil
	}

	array := make([]RESP, length)

	for i := 0; i < length; i++ {
		val, err := p.Parse()
		if err != nil {
			return nil, err
		}
		array[i] = *val
	}

	return &RESP{Type: "array", Array: array}, nil

}

func (p *RedisParser) parseInteger() (*RESP, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	num, err := strconv.Atoi(line)

	if err != nil {
		return nil, fmt.Errorf("invalid integer")
	}

	return &RESP{Type: "integer", Number: num}, nil
}

func (p *RedisParser) parseBulkString() (*RESP, error) {
	line, err := p.readLine()

	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(line)

	if err != nil {
		return nil, fmt.Errorf("invalid bulk string, expected to get length from bulk string")
	}

	if length == -1 {
		return &RESP{Type: "null"}, nil
	}

	if length == 0 {
		_, err := p.readLine()

		if err != nil {
			return nil, err
		}
		return &RESP{Type: "string", String: ""}, nil
	}

	buf := make([]byte, length)

	_, err = io.ReadFull(p.reader, buf)

	if err != nil {
		return nil, err
	}

	_, err = p.readLine()
	if err != nil {
		return nil, err
	}

	return &RESP{Type: "string", String: string(buf)}, nil
}

/*reads a line until it reaches \r\n */
func (p *RedisParser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	if len(line) < 2 || line[len(line)-2] != '\r' {
		return "", fmt.Errorf("invalid line ending")
	}

	return line[:len(line)-2], nil
}

func (v *RESP) convertToCommand() []string {
	if v.Type != "array" {
		return []string{v.String}
	}

	cmd := make([]string, len(v.Array))

	for i, arg := range v.Array {
		cmd[i] = arg.String
	}

	return cmd
}

func FormatResp(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("+%s\r\n", v)
	case int:
		return fmt.Sprintf(":%d\r\n", v)
	case error:
		return fmt.Sprintf("-%s\r\n", v.Error())
	case []string:
		finalBulkStr := make([]string, len(v))
		finalBulkStr = append(finalBulkStr, fmt.Sprintf("*%d\r\n", len(v)))

		for i := 0; i < len(v); i++ {
			appendStr := fmt.Sprintf("$%d\r\n%s\r\n", len(v[i]), v[i])
			finalBulkStr = append(finalBulkStr, appendStr)
		}

		res := strings.Join(finalBulkStr, "")
		fmt.Printf("result:: %s", res)
		return res
	case nil:
		return "$-1\r\n"
	default:
		return fmt.Sprintf("+%v\r\n", v)
	}
}

func handleCommand(cmd []string) string {
	if len(cmd) == 0 {
		return FormatResp(errors.New("-ERR: empty command"))
	}

	command := strings.ToUpper(cmd[0])

	switch command {
	case "PING":
		if len(cmd) == 1 {
			return FormatResp("PONG")
		}

		return FormatResp(cmd[1])

	case "ECHO":
		if len(cmd) < 2 {
			return FormatResp(errors.New("-ERR ECHO needs arguments"))
		}

		return FormatResp(cmd[1])

	case "SET":
		if len(cmd) < 3 {
			return FormatResp(errors.New("-ERR wrong number of argumnts for 'SET' command"))
		}

		if err := setValue(cmd[1:]); err != nil {
			return FormatResp(errors.New("-ERR something went wrong"))
		}

		return FormatResp("OK")

	case "GET":
		if len(cmd) < 2 {
			return FormatResp(errors.New("-ERR wrong number of argumnts for 'GET' command"))
		}

		val, err := getValue(cmd[1])

		if err != nil {
			return FormatResp(errors.New("1"))
		}

		if val != nil {
			return FormatResp(val)
		}

		return FormatResp(nil)
	case "CONFIG":
		if len(cmd) < 3 {
			return FormatResp(errors.New("-ERR wrong number of argumnts for 'CONFIG' command, required 3 args"))
		}

		if strings.Compare(cmd[1], "GET") != 0 {
			return FormatResp(errors.New("-ERR expected to receive 'GET' as second argument"))
		}

		rdbStorageConfig := getStorageConfig(cmd[2])

		return FormatResp(rdbStorageConfig)

	default:
		return FormatResp(fmt.Errorf("-ERR unknown command '%s'", command))
	}
}
