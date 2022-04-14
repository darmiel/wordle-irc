package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

var commands []*cli.Command

func main() {
	// bootstrap: CLI for bootstrapping the ic
	bootstrap := &cli.App{
		Name:                 "WordleIIRC",
		Version:              "1.0.0",
		Description:          "Play Wordle in IRC!",
		Commands:             commands,
		Flags:                nil,
		EnableBashCompletion: true,
		Action:               nil,
		Authors: []*cli.Author{{
			Name:  "darmiel",
			Email: "hi@d2a.io",
		}},
		UseShortOptionHandling: true,
	}
	if err := bootstrap.Run(os.Args); err != nil {
		panic(err)
	}
}
