package main

import (
	"fmt"
	"github.com/apex/log"
	ch "github.com/apex/log/handlers/cli"
	"github.com/urfave/cli/v2"
	"goircwordle/internal/ic"
	"net"
)

func init() {
	log.SetHandler(ch.Default)
	var (
		host string
		port uint

		pass string
		name string
	)

	commands = append(commands, &cli.Command{
		Name:        "connect",
		Aliases:     []string{"c"},
		Usage:       "",
		Description: "Connect to an IRC server",
		Action: func(ctx *cli.Context) (err error) {
			addr := fmt.Sprintf("%s:%d", host, port)
			log.Infof("Connecting to %s", addr)

			var conn net.Conn
			if conn, err = net.Dial("tcp", addr); err != nil {
				return
			}

			client := ic.NewClient(conn, name, pass)
			return client.Run()
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       "irc.d2a.io",
				Destination: &host,
			},
			&cli.UintFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       6667,
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "pass",
				Aliases:     []string{"x"},
				Destination: &pass,
			},
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"u", "n"},
				Destination: &name,
				Value:       "wordlebot",
			},
		},
	})
}
