package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	if len(args) > 4 || len(args) < 2 || len(args) == 3 {
		return "", errors.New("invalid amount of arguments")
	}
	options, err := extractOptions(args)
	if err != nil {
		return "", err
	}
	if len(args) == 2 {
		h.storage.Set(options[0], options[1])
	} else {
		ms, err := strconv.Atoi(options[3])
		if err != nil {
			return "", err
		}
		h.storage.SetWithExpiration(options[0], options[1], time.Millisecond*time.Duration(ms))
	}

	return Encode(Result{
		Type:  RedisSimpleString,
		Value: "OK",
	})
}

func extractOptions(args []Result) ([]string, error) {
	var options []string
	key, ok := args[0].Value.(string)
	if !ok {
		return nil, errors.New("failed to convert key to string")
	}
	value, ok := args[1].Value.(string)
	if !ok {
		return nil, errors.New("failed to convert value to string")
	}
	options = append(options, key, value)
	if len(args) > 2 {
		px, ok := args[2].Value.(string)
		if !ok {
			return nil, errors.New("failed to convert key to string")
		}
		if strings.ToLower(px) != "px" {
			return nil, errors.New("expected px to be the option")
		}
		duration, ok := args[3].Value.(string)
		if !ok {
			return nil, errors.New("failed to convert value to string")
		}
		options = append(options, px, duration)
	}
	return options, nil
}
