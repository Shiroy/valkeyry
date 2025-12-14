package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/client"
)

type Echo struct{}

func (e Echo) Mnemonic() string {
	return "ECHO"
}

func NewEcho() *Echo {
	return &Echo{}
}

func (e Echo) Handle(c *client.Session, command []string) error {
	if len(command) < 2 {
		return fmt.Errorf("invalid command. Not enough arguments")
	}

	err := c.SendString(command[1])
	return err
}
