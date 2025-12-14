package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Session struct {
	socket net.Conn
	reader *bufio.Reader
}

func NewSession(socket net.Conn) *Session {
	return &Session{
		socket: socket,
		reader: bufio.NewReader(socket),
	}
}

func (c *Session) ParseCommand() ([]string, error) {
	val, err := c.reader.Peek(1)

	if err != nil {
		return nil, err
	}

	switch val[0] {
	case '*':
		cmd, err := c.parseArray()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	default:
		return nil, fmt.Errorf("unsupported command. Received %v", val[0])
	}
}

func (c *Session) parseArray() ([]string, error) {
	val, err := c.reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if val != '*' {
		return nil, fmt.Errorf("invalid starting byte for array. Got '%q', expected '*'", val)
	}

	size, err := c.parseInteger()
	if err != nil {
		return nil, err
	}

	err = c.parseSeparator()
	if err != nil {
		return nil, err
	}

	array := make([]string, 0, size)

	for i := 0; i < size; i++ {
		next, err := c.reader.Peek(1)
		if err != nil {
			return nil, err
		}

		switch next[0] {
		case '$':
			value, err := c.parseString()
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

func (c *Session) parseString() (value string, err error) {
	// Read '$'
	val, err := c.reader.ReadByte()
	if err != nil {
		return "", err
	}

	if val != '$' {
		return "", fmt.Errorf("invalid starting byte for string. Got '%q', expected '$'", val)
	}

	size, err := c.parseInteger()
	if err != nil {
		return "", err
	}

	err = c.parseSeparator()
	if err != nil {
		return "", err
	}

	buf := make([]byte, 0, size)

	for i := 0; i < size; i++ {
		c, err := c.reader.ReadByte()
		if err != nil {
			return "", err
		}

		buf = append(buf, c)
	}

	err = c.parseSeparator()
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (c *Session) parseInteger() (int, error) {
	buf := make([]byte, 0, 64)

	for {
		next, err := c.reader.Peek(1)
		if err != nil {
			return 0, err
		}

		if next[0] < '0' || next[0] > '9' {
			break
		}

		buf = append(buf, next[0])

		_, err = c.reader.ReadByte()
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

func (c *Session) parseSeparator() error {
	first, err := c.reader.ReadByte()
	if err != nil {
		return err
	}

	second, err := c.reader.ReadByte()
	if err != nil {
		return err
	}

	if first != '\r' || second != '\n' {
		return errors.New("invalid separator")
	}

	return nil
}

func (c *Session) SendSimpleString(value string) error {
	_, err := c.socket.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
	return err
}

func (c *Session) SendString(value string) error {
	_, err := c.socket.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))

	return err
}

func (c *Session) SendPong() error {
	return c.SendSimpleString("PONG")
}

func (c *Session) SendErrorString(msg string) error {
	_, err := c.socket.Write([]byte(fmt.Sprintf("-%s\r\n", msg)))
	return err
}

func (c *Session) SendNullBulkString() error {
	_, err := c.socket.Write([]byte("$-1\r\n"))
	return err
}
