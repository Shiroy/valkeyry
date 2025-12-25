package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	command []string
	pos     int
}

var EndOfSeq = errors.New("end of sequence")

func NewParser(command []string) Parser {
	return Parser{
		command: command,
		pos:     0,
	}
}

func (p *Parser) Peek() (string, error) {
	if p.pos >= len(p.command) {
		return "", EndOfSeq
	} else {
		return p.command[p.pos], nil
	}
}

func (p *Parser) Read() (string, error) {
	if p.pos >= len(p.command) {
		return "", EndOfSeq
	} else {
		value := p.command[p.pos]
		p.pos += 1
		return value, nil
	}
}

func (p *Parser) ReadInt() (int, error) {
	val, err := p.Read()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func (p *Parser) ReadLiteral(literal string, caseInsensitive bool) error {
	val, err := p.Read()
	if err != nil {
		return err
	}

	if caseInsensitive {
		literal = strings.ToUpper(literal)
		val = strings.ToUpper(val)
	}

	if val != literal {
		return fmt.Errorf("literal mismatch: expected '%s' found '%s'", literal, val)
	} else {
		return nil
	}
}
