package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	redisdir        = flag.String("dir", "", "Path to the directory where the rbd file is stored ")
	redisdbfilename = flag.String("dbfilename", "", "name of the rdb file (example: rdbfile)")
)

type RdbStorage struct {
	dir      string
	filename string
}

type KeyValue struct {
	Data []string
	Time *time.Time
}

var DB sync.Map

// redis-cli SET foo bar px 100 - # Sets the key "foo" to "bar" with an expiry of 100 milliseconds
func setValue(keyval []string) error {
	if len(keyval) < 2 {
		return fmt.Errorf("did not find key and value")
	}

	data := KeyValue{}

	if len(keyval) > 2 {
		data.Data = keyval[1 : len(keyval)-2]
		exp, err := strconv.Atoi(keyval[len(keyval)-1])
		if err != nil {
			return fmt.Errorf("could not convert time to int")
		}

		durationToLast := time.Now().Add(time.Duration(exp) * time.Millisecond)
		data.Time = &durationToLast
	} else {
		data.Data = keyval[1:]
		data.Time = nil
	}

	DB.Store(keyval[0], data)

	return nil
}

func getValue(key string) (interface{}, error) {
	val, ok := DB.Load(key)

	if !ok {
		return "", fmt.Errorf("value not found")
	}

	dataKeyVal, ok := val.(KeyValue)

	if dataKeyVal.Time == nil {
		return strings.Join(dataKeyVal.Data, " "), nil
	}

	if time.Now().After(*dataKeyVal.Time) {
		return nil, nil
	}

	return strings.Join(dataKeyVal.Data, " "), nil
}

/*
 Set storage variables and parse them
*/

func rdbFileConfig() {
	flag.Parse()
}

func getStorageConfig(config string) []string {
	switch config {
	case "dir":
		return []string{"dir", *redisdir}
	case "dbfilename":
		return []string{"dbfilename", *redisdbfilename}
	default:
		return nil
	}
}
