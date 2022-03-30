package main

import "strings"

type Channel string

func findGameInChannel(channel Channel) *Game {
	c := channel.Normalize()
	gs, ok := games[c]
	if ok {
		for _, game := range gs {
			if game.channel == c && game.active {
				return game
			}
		}
	}
	return nil
}

func (c Channel) Normalize() string {
	channel := string(c)
	for strings.HasPrefix(channel, "#") {
		channel = channel[1:]
	}
	return strings.TrimSpace(strings.ToLower(channel))
}

func (c Channel) String() string {
	return "#" + c.Normalize()
}
