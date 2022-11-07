package main

import (
	"errors"
	"fmt"
	"strings"
)

type Handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Handle(command []byte) (string, error) {
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

	return h.handleCommand(cmd, args)

}

func (h *Handler) handleCommand(cmd Result, args []Result) (string, error) {
	switch cmd.Value {
	case "ping":
		return handlePing()
	case "echo":
		return handleEcho(args)
	case "get":
		return h.handleGet(args)
	case "set":
		return h.handleSet(args)
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

func handleEcho(args []Result) (string, error) {
	var stringSlice []string
	for _, result := range args {
		value, ok := result.Value.(string)
		if !ok {
			return "", errors.New("failed to convert value to string")
		}
		stringSlice = append(stringSlice, value)
	}

	return Encode(Result{
		Type:  RedisBulkString,
		Value: strings.Join(stringSlice, " "),
	})
}

func (h *Handler) handleGet(args []Result) (string, error) {
	if len(args) > 1 {
		return "", errors.New("invalid amount of arguments")
	}
	key, ok := args[0].Value.(string)
	if !ok {
		return "", errors.New("failed to convert key to string")
	}
	value, ok := h.storage.Get(key)
	response := Result{
		Type: RedisBulkString,
	}
	if ok {
		response.Value = value
		return Encode(response)
	}
	return Encode(response)

}

func (h *Handler) handleSet(args []Result) (string, error) {
	if len(args) > 2 {
		return "", errors.New("invalid amount of arguments")
	}
	key, ok := args[0].Value.(string)
	if !ok {
		return "", errors.New("failed to convert key to string")
	}
	value, ok := args[0].Value.(string)
	if !ok {
		return "", errors.New("failed to convert value to string")
	}

	h.storage.Set(key, value)

	return Encode(Result{
		Type:  RedisSimpleString,
		Value: "OK",
	})
}
