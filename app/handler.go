package main

import (
	"errors"
	"fmt"
)

func Handle(command []byte) (string, error) {
	result := Decode(command)
	if result.Type != RedisArray {
		return "", errors.New("command can ony be encoded as Redis Array")
	}
	value, ok := result.Value.([]Result)
	if !ok {
		return "", errors.New("can't convert value to slice")
	}

	cmd := value[0]
	args := value[1:]

	return handleCommand(cmd, args)

}

func handleCommand(cmd Result, args []Result) (string, error) {
	switch cmd.Value {
	case "ping":
		return handlePing()
	case "echo":
		return "", nil
	default:
		fmt.Println(cmd.Value)
		return "", errors.New("unrecognized command: ")
	}

}

func handlePing() (string, error) {
	return Encode(Result{
		Type:  RedisSimpleString,
		Value: "PONG",
	})
}
