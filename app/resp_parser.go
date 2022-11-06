package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type RedisType int

const (
	RedisSimpleString RedisType = iota
	RedisError
	RedisInt
	RedisBulkString
	RedisArray
	RedisNull
)

const crlfLength = 4

type Result struct {
	Type  RedisType
	Value interface{}
}

func Decode(data []byte) Result {
	v, _ := decode(data)
	return v
}

func Encode(value Result) (string, error) {
	switch value.Type {
	case RedisSimpleString:
		return fmt.Sprintf("+%s\r\n", value.Value), nil
	case RedisError:
		return fmt.Sprintf("-%s\r\n", value.Value), nil
	case RedisInt:
		return fmt.Sprintf(":%d\r\n", value.Value), nil
	case RedisBulkString:
		if value.Value == nil {
			return "$-1\r\n", nil
		}
		v, ok := value.Value.(string)
		if !ok {
			return "", errors.New("cannot get bulk string value")
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v), nil
	default:
		return "", nil
	}
}

func decodeString(value []byte) (Result, int) {
	i, end := matchCRLF(value)
	return Result{
		Type:  RedisSimpleString,
		Value: string(value[1:i]),
	}, end
}

func decodeError(value []byte) (Result, int) {
	i, end := matchCRLF(value)
	return Result{
		Type:  RedisError,
		Value: string(value[1:i]),
	}, end
}

func decodeInt(value []byte) (Result, int) {
	i, end := matchCRLF(value)
	v, err := strconv.Atoi(string(value[1:i]))
	if err != nil {
		return Result{
			Type:  RedisError,
			Value: fmt.Sprintf("cannot convert %s to int", string(value[1:i])),
		}, end
	}
	return Result{
		Type:  RedisInt,
		Value: v,
	}, end
}

func decodeBulkString(value []byte) (Result, int) {
	i, end := matchCRLF(value)
	size, err := strconv.Atoi(string(value[1:i]))
	if err != nil {
		return Result{
			Type:  RedisError,
			Value: "invalid bulk string size",
		}, end
	}

	if size == -1 {
		return Result{
			Type: RedisNull,
		}, end
	}

	firstEnd := end
	valueAfterSize := value[end:]
	i, end = matchCRLF(valueAfterSize)
	v := string(valueAfterSize[:i])
	return Result{
		Type:  RedisBulkString,
		Value: v,
	}, firstEnd + end
}

func decodeArray(value []byte) (Result, int) {
	i, end := matchCRLF(value)
	size, err := strconv.Atoi(string(value[1:i]))
	if err != nil {
		return Result{
			Type:  RedisError,
			Value: "invalid array size",
		}, end
	}

	if size == -1 {
		return Result{
			Type: RedisNull,
		}, end
	}

	if size == 0 {
		return Result{
			Type: RedisArray,
		}, end
	}

	vals := make([]Result, 0, size)
	currentPos := end
	for i := 0; i < size; i++ {
		v, p := decode(value[currentPos:])
		vals = append(vals, v)
		currentPos += p
	}

	return Result{
		Type:  RedisArray,
		Value: vals,
	}, currentPos

}

func decode(value []byte) (Result, int) {
	dataType := string(value[:1])
	switch dataType {
	case "+":
		return decodeString(value)
	case "-":
		return decodeError(value)
	case ":":
		return decodeInt(value)
	case "$":
		return decodeBulkString(value)
	case "*":
		return decodeArray(value)
	default:
		return Result{
			Type:  RedisError,
			Value: fmt.Sprintf("unrecognzied data type %s", dataType),
		}, 0
	}
}

func matchCRLF(value []byte) (pos, end int) {
	i := bytes.IndexByte(value, '\\')
	return i, i + crlfLength
}
