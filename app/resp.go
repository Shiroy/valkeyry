package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func ParseResp(reader *bufio.Reader) ([]string, error) {
	val, err := reader.Peek(1)

	if err != nil {
		return nil, err
	}

	switch val[0] {
	case '*':
		cmd, err := ParseArray(reader)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		return nil, fmt.Errorf("unsupported command. Received %v", val[0])
	}
}

func ParseArray(reader *bufio.Reader) ([]string, error) {
	val, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if val != '*' {
		return nil, fmt.Errorf("invalid starting byte for array. Got '%q', expected '*'", val)
	}

	size, err := ParseInteger(reader)
	if err != nil {
		return nil, err
	}

	fmt.Println("Array size:", size)
	err = ParseSeparator(reader)
	if err != nil {
		return nil, err
	}
	fmt.Println("Received array separator")

	array := make([]string, 0, size)

	for i := 0; i < size; i++ {
		next, err := reader.Peek(1)
		if err != nil {
			return nil, err
		}

		switch next[0] {
		case '$':
			value, err := ParseString(reader)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
			break
		default:
			return nil, fmt.Errorf("unsupported element type: %q", next)
		}
	}

	return array, nil
}

func ParseString(reader *bufio.Reader) (value string, err error) {
	// Read '$'
	val, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	if val != '$' {
		return "", fmt.Errorf("invalid starting byte for string. Got '%q', expected '$'", val)
	}

	size, err := ParseInteger(reader)
	if err != nil {
		return "", err
	}

	err = ParseSeparator(reader)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 0, size)

	for i := 0; i < size; i++ {
		c, err := reader.ReadByte()
		if err != nil {
			return "", err
		}

		buf = append(buf, c)
	}

	err = ParseSeparator(reader)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func ParseInteger(reader *bufio.Reader) (int, error) {
	buf := make([]byte, 0, 64)

	for {
		next, err := reader.Peek(1)
		if err != nil {
			return 0, err
		}

		if next[0] < '0' || next[0] > '9' {
			break
		}

		buf = append(buf, next[0])

		_, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	}

	val, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}

	return val, nil
}

func ParseSeparator(reader *bufio.Reader) error {
	first, err := reader.ReadByte()
	if err != nil {
		return err
	}

	second, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if first != '\r' || second != '\n' {
		return errors.New("invalid separator")
	}

	return nil
}
