package main

import (
	"fmt"
	"gopkg.in/irc.v3"
	"log"
	"net"
	"strconv"
	"strings"
)

type GameState struct {
	Word  string
	Tries int
}

var States = make(map[string]*GameState)

func start(channel, word string, client *irc.Client) {
	States[channel] = &GameState{
		Word:  strings.ToUpper(word),
		Tries: 0,
	}

	// print hello world message
	send := func(msg string) {
		_ = client.WriteMessage(&irc.Message{
			Command: "PRIVMSG",
			Params: []string{
				channel,
				msg,
			},
		})
	}

	send("Hello 👋 I am a Wordle bot for the cool Qwiri IRC.")
	send(fmt.Sprintf("The current word is %s%d%s characters long.",
		Color("11"), len(word), Color("")))
	send(fmt.Sprintf("Use %sguess <WORD>%s to guess a word.",
		Color("11"), Color("")))
	send("(?) " + strings.Repeat("_ ", len(word)))
}

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

			send := func(msg string) {
				_ = c.WriteMessage(&irc.Message{
					Command: "PRIVMSG",
					Params: []string{
						m.Params[0],
						msg,
					},
				})
			}

			switch m.Command {
			case "001":
				if err := c.Write("JOIN #wordle"); err != nil {
					fmt.Println("Cannot join channel:", err)
				} else {
					fmt.Println("Joined wordle channel!")
				}

			case "JOIN":
				log.Println("Joined Channel", m.Params[0])
				States[m.Params[0]] = &GameState{}

			case "PRIVMSG":
				var message = m.Trailing()

				if !c.FromChannel(m) {
					fmt.Printf("Got private message: %+v\n", m)

					// set new word
					if !strings.HasPrefix(message, "set ") {
						return
					}
					split := strings.Split(message, " ")
					if len(split) != 3 {
						log.Println("Split was not 3")
						return
					}

					channel := split[1]
					word := split[2]

					if !checkWordValid(word) {
						return
					}

					_, ok := States[channel]
					if !ok {
						_ = c.Write("JOIN " + channel)
						return
					}

					start(channel, word, c)
					return
				}

				// check if we have a state for the channel
				state, ok := States[m.Params[0]]
				if !ok || state.Word == "" {
					return
				}

				if !strings.HasPrefix(message, "guess ") {
					return
				}

				guess := strings.ToUpper(strings.TrimSpace(message[len("guess "):]))
				if guess == "" {
					return
				}

				if len(guess) != len(state.Word) {
					send("🤬 This guess is too long/short")
					return
				}

				state.Tries++

				var bob strings.Builder
				bob.WriteString("(" + strconv.Itoa(state.Tries) + ") ")

				correct := 0
				for i, cg := range guess {
					cc := state.Word[i]

					var color string
					if uint8(cg) == cc { // correct word
						correct++
						color = "0,3"
					} else if InWord(state.Word, cg) {
						color = "0,7"
					} else {
						color = "0,14"
					}

					bob.WriteRune(3)
					bob.WriteString(color)
					bob.WriteString(fmt.Sprintf(" %c ", cg))
					bob.WriteRune(3)
					bob.WriteRune(' ')
				}

				send(bob.String())

				if correct == len(state.Word) {
					send(fmt.Sprintf("✅ Nice! You got the word in %d tries.", state.Tries))
					state.Word = ""
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

func InWord(word string, c int32) bool {
	for _, cc := range word {
		if c == cc {
			return true
		}
	}
	return false
}

// checkWordValid checks if the given word is a Heterogram
func checkWordValid(word string) bool {
	for _, char := range word {
		if strings.Count(word, string(char)) > 1 {
			return false
		}
	}
	return true
}

func Color(a string) string {
	return string(rune(3)) + a
}
