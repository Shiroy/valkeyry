package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/client"
)

type Ping struct{}

func (p Ping) Mnemonic() string {
	return "PING"
}

func NewPing() *Ping {
	return &Ping{}
}

func (p Ping) Handle(c *client.Session, _ []string) error {
	err := c.SendPong()
	return err
}
