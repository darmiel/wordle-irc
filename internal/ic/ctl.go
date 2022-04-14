package ic

import "github.com/urfave/cli/v2"

func (c *Client) ctlApp() *cli.App {
	return &cli.App{
		Name:    "wordlectl",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:        "start",
				Aliases:     []string{"s"},
				Description: "Starts a new wordle session",
				Category:    "control",
				Action: func(ctx *cli.Context) error {
					// TODO: Implement me
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "channel",
						Aliases:  []string{"c"},
						Usage:    "The channel to start the wordle session",
						Required: true,
					},
				},
			},
		},
		Flags:                  nil,
		UseShortOptionHandling: true,
	}
}
