package main

import (
	"fmt"
	"gopkg.in/irc.v3"
	"log"
	"net"
	"strings"
)

var games = make(map[string][]*Game)

func main() {
	conn, err := net.Dial("tcp", "irc.d2a.io:6667")
	if err != nil {
		log.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick: "wordlebot",
		Pass: "ABAP",
		User: "wordlebot",
		Name: "Wordle",
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {

			switch m.Command {
			case "001":
				log.Println("001")
				if err := start("wordle", "FUCHS", c); err != nil {
					fmt.Println("Cannot start wordle test:", err)
				}

			case "JOIN":
				log.Println("JOIN")

			case "PRIVMSG":
				var (
					message = m.Trailing()
					channel = Channel(m.Params[0])
					game    = findGameInChannel(channel)
				)

				log.Println("PRIVMSG", channel.String(), "::", message)

				// no game found in channel
				if game == nil {
					return
				}

				if !strings.HasPrefix(message, "guess ") &&
					!strings.HasPrefix(message, "g ") {
					return
				}

				guess := Word(strings.SplitN(message, " ", 2)[1]).Normalize()
				if guess == "" {
					return
				}

				if err := game.handleGuess(guess, m.User); err != nil {
					log.Println("ERR | Handling guess in channel", channel, "::", err)
					return
				}
			}
		}),
	}

	// Create the client
	client := irc.NewClient(conn, config)
	err = client.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func start(channel Channel, word Word, client *irc.Client) error {
	c := channel.Normalize()
	game := &Game{
		word:    word.Normalize(),
		channel: c,
		tries:   make(map[string]uint),
		client:  client,
		active:  true,
	}

	// append game
	gs, _ := games[c]
	gs = append(gs, game)
	games[c] = gs

	return game.hello()
}
