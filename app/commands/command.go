package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/client"
)

type Command interface {
	Handle(c *client.Session, command []string) error
	Mnemonic() string
}
