package main

import (
	"errors"
	"fmt"
	"gopkg.in/irc.v3"
	"strings"
)

var ErrGameInactive = errors.New("game is inactive")

type Game struct {
	word    Word
	channel Channel
	tries   map[string]uint
	client  *irc.Client
	active  bool

	guessed []rune
	hard    bool
}

func (g *Game) addTry(user string) uint {
	t, _ := g.tries[user]
	t++
	g.tries[user] = t
	return t
}

func (g *Game) sumTries() uint {
	var res uint
	for _, t := range g.tries {
		res += t
	}
	return res
}

func (g *Game) sendChannelMessage(msg string) error {
	return g.channel.Send(g.client, msg)
}

func (g *Game) reply(user string, msg ...string) (err error) {
	for _, m := range msg {
		if err = g.send(user + ": " + m); err != nil {
			return
		}
	}
	return nil
}

func (g *Game) send(msg ...string) error {
	for _, m := range msg {
		if err := g.sendChannelMessage(m); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) hello() (err error) {
	// join channel
	if err = g.client.Write("JOIN " + string(g.channel)); err != nil {
		return
	}
	// print hello
	return g.send(
		"Hello 👋 Let's play WordleIIRC!",
		fmt.Sprintf(
			"The current word is %s%d%s characters long.",
			ColorCyan.String(), len(g.word), StyleReset.String(),
		),
		fmt.Sprintf(
			"Use %sg <word>%s to guess a word.",
			ColorCyan.String(), StyleReset.String(),
		),
		fmt.Sprintf(
			"The word is: %s",
			g.word.Print(Word(strings.Repeat("_", len(g.word)))),
		),
	)
}

func (g *Game) handleGuess(guess Word, user string) error {
	if !g.active {
		return ErrGameInactive
	}

	// check if length's matches
	if len(g.word) != len(guess) {
		return g.reply(user, fmt.Sprintf(
			"%s :: You're entered a %d (req %d) char long word.",
			ColorRedBG.Enclose("ERR"), len(guess), len(g.word),
		))
	}

	// check hard mode
	if g.hard {
		for i, gu := range guess {
			co := g.guessed[i]
			if co == 0 { // not guessed
				continue
			}
			if gu != co {
				return g.reply(user, fmt.Sprintf(
					"%s :: You're playing %shard mode%s. This word doesn't match your guesses.",
					ColorRedBG.Enclose("ERR"), ColorCyan.String(), StyleReset.String(),
				))
			}
		}
	}

	// save correct guesses for hard mode
	for i := range guess {
		gu := guess.At(i)
		if g.word.At(i) == gu {
			g.guessed[i] = gu
		}
	}

	// add try for user
	tries := g.addTry(user)
	sumTries := g.sumTries()

	if err := g.sendChannelMessage(fmt.Sprintf(
		"%s :: %s",
		ColorCyanBG.Enclose(fmt.Sprintf("🔰 [u:%d;s:%d]", tries, sumTries)),
		g.word.Print(guess),
	)); err != nil {
		return err
	}

	// check if the guessed word is correct
	if guess == g.word {
		g.active = false // set game as incorrect to prevent any other guesses

		if err := g.sendChannelMessage(fmt.Sprintf(
			"%s :: The word was guessed after a total of %d tries.",
			ColorGreenBG.Enclose("✅ NICE"), sumTries,
		)); err != nil {
			return err
		}
	}

	return nil
}
