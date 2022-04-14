package ic

import (
	"github.com/apex/log"
	"github.com/urfave/cli/v2"
	"goircwordle/internal/common"
	"gopkg.in/irc.v3"
	"net"
	"strings"
)

type Client struct {
	client *irc.Client
	games  map[common.Channel]*Game
	ctl    *cli.App
}

func NewClient(conn net.Conn, name, pass string) (c *Client) {
	c = &Client{
		games: make(map[common.Channel]*Game),
	}
	config := irc.ClientConfig{
		Nick:    name,
		Pass:    pass,
		User:    name,
		Name:    name,
		Handler: irc.HandlerFunc(c.Handler),
	}
	c.client = irc.NewClient(conn, config)
	c.ctl = c.ctlApp()
	return
}

func (c *Client) Run() error {
	return c.client.Run()
}

func (c *Client) Handler(ic *irc.Client, im *irc.Message) {
	if im.Command == "001" {
		_ = c.start(common.ChannelOf("wordle"), common.WordOf("BROT"))
		return
	}

	if im.Command != "PRIVMSG" {
		return
	}
	var (
		message = im.Trailing()
		channel = common.ChannelOf(im.Params[0])
	)

	// private message: admin command mode
	if !ic.FromChannel(im) {
		if err := c.ctl.Run(strings.Split(message, " ")); err != nil {
			log.WithError(err).Warn("cannot run private message parsing")
			return
		}
	}

	// no game found in channel
	game, ok := c.games[channel]
	if !ok || game == nil {
		return
	}

	if !strings.HasPrefix(message, "guess ") &&
		!strings.HasPrefix(message, "g ") {
		return
	}

	guess := common.WordOf(strings.SplitN(message, " ", 2)[1])
	if guess.Word() == "" {
		return
	}

	if err := game.handleGuess(guess, im.User); err != nil {
		log.WithError(err).
			WithField("channel", channel).
			Warn("handling guess failed")
		return
	}
}

func (c *Client) start(channel common.Channel, word common.Word) error {
	log.Infof("New word '%s' in channel %s", strings.ToUpper(word.Word()), channel)
	game := &Game{
		word:    word,
		channel: channel,
		tries:   make(map[string]uint),
		client:  c.client,
		active:  true,
		hard:    true,
		guessed: make([]rune, len(word.Word())),
	}
	c.games[channel] = game
	return game.hello()
}
