package common

import (
	"gopkg.in/irc.v3"
	"strings"
)

// Channel wrapper
type Channel interface {
	Channel() string
	Send(client *irc.Client, msg string) error
}

// channel implements Channel
type channel string

func (c channel) Channel() string {
	return string(c)
}

func (c channel) Send(client *irc.Client, msg string) error {
	return client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params: []string{
			string(c),
			msg,
		},
	})
}

///

func ChannelOf(c string) Channel {
	return channel("#" + strings.TrimLeft(strings.TrimSpace(strings.ToLower(c)), "#"))
}
