package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KeyValue struct {
	Data []string
	Time time.Time
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

		durationToLast := time.Duration(exp) * time.Millisecond
		data.Time = time.Now().Add(durationToLast)
	} else {
		data.Data = keyval[1:]
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

	if time.Now().After(dataKeyVal.Time) {
		return "", fmt.Errorf("value not found")
	}

	return strings.Join(dataKeyVal.Data, " "), nil
}
