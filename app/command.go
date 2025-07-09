package main

import (
	"fmt"
	"sync"
)

var DB sync.Map

func setValue(keyval []string) error {
	if len(keyval) < 2 {
		return fmt.Errorf("did not find key and value")
	}

	DB.Store(keyval[0], keyval[1:])

	return nil
}

func getValue(key string) (interface{}, error) {
	val, ok := DB.Load(key)

	if !ok {
		return "", fmt.Errorf("value not found")
	}

	return val, nil
}
