package main

import (
	"gopkg.in/irc.v3"
	"strings"
)

type Channel string

func channelOf(channel string) Channel {
	return Channel("#" + strings.TrimLeft(strings.TrimSpace(strings.ToLower(channel)), "#"))
}

func findGameInChannel(channel Channel) *Game {
	gs, ok := games[channel]
	if ok {
		for _, game := range gs {
			if game.channel == channel && game.active {
				return game
			}
		}
	}
	return nil
}

func (c Channel) Send(client *irc.Client, msg string) error {
	return client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params: []string{
			string(c),
			msg,
		},
	})
}
