package main

import (
	"fmt"
)

var store = make(map[string]string)

func setValue(keyval []string) error {
	if len(keyval) != 2 {
		return fmt.Errorf("did not find key and value")
	}

	store[keyval[0]] = keyval[1]

	return nil
}

func getValue(key string) (string, error) {
	val, found := store[key]

	if !found {
		return "", fmt.Errorf("value not found")
	}

	return val, nil
}
