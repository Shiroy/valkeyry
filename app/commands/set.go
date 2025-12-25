package commands

import (
	"errors"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/memory/values"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SetParams struct {
	fx.In
	Cache *memory.Cache
	Log   *zap.Logger
}

type Set struct {
	cache *memory.Cache
	log   *zap.Logger
}

func (s *Set) Handle(c *client.Session, command []string) error {
	input, err := SetInputFromCommand(command)

	if err != nil {
		return err
	}

	if input.expireAt != nil {
		s.cache.SetWithExpiration(input.key, values.NewValueString(input.value), *input.expireAt)
	} else {
		s.cache.Set(input.key, values.NewValueString(input.value))
	}

	return c.SendSimpleString("OK")
}

func (s *Set) Mnemonic() string {
	return "SET"
}

type SetInput struct {
	key      string
	value    string
	expireAt *int64 // expiration timestamp in milliseconds
}

func SetInputFromCommand(command []string) (*SetInput, error) {
	parser := NewParser(command)

	var expireAt *int64 = nil

	if err := parser.ReadLiteral("SET", true); err != nil {
		return nil, err
	}

	key, err := parser.Read()
	if err != nil {
		return nil, err
	}

	value, err := parser.Read()
	if err != nil {
		return nil, err
	}

	next, err := parser.Peek()
	if errors.Is(err, EndOfSeq) {
		// No-op
	} else if err != nil {
		return nil, err
	} else {
		switch strings.ToUpper(next) {
		case "EX":
			expire, err := ParseEx(parser)
			if err != nil {
				return nil, err
			}
			expireAt = &expire
		case "PX":
			expire, err := ParseNx(parser)
			if err != nil {
				return nil, err
			}
			expireAt = &expire
		}
	}

	return &SetInput{
		key:      key,
		value:    value,
		expireAt: expireAt,
	}, nil
}

func ParseNx(parser Parser) (int64, error) {
	if err := parser.ReadLiteral("PX", true); err != nil {
		return 0, err
	}

	expirationInMS, err := parser.ReadInt()
	if err != nil {
		return 0, err
	}

	expire := time.Now().UnixMilli() + int64(expirationInMS)
	return expire, nil
}

func ParseEx(parser Parser) (int64, error) {
	if err := parser.ReadLiteral("EX", true); err != nil {
		return 0, err
	}

	expirationInSecond, err := parser.ReadInt()
	if err != nil {
		return 0, err
	}

	expire := (time.Now().Unix() + int64(expirationInSecond)) * 1000

	return expire, nil
}

func NewSet(p SetParams) *Set {
	return &Set{
		log:   p.Log.With(zap.String("command", "SET")),
		cache: p.Cache,
	}
}
