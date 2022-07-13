package main

import (
	"os"
	"strconv"

	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Serve",
		Usage:       "The Serve",
		Description: "Server static files",
		Version:     Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Value:   "9000",
				Usage:   "The port to listen on",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "prefix",
				Value:   "/",
				Usage:   "The prefix to listen on",
				Aliases: []string{},
			},
			&cli.StringFlag{
				Name:    "dir",
				Value:   fs.CurrentDir(),
				Usage:   "The dir to listen on",
				Aliases: []string{"d"},
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			if os.Getenv("PORT") != "" {
				port = os.Getenv("PORT")
			}

			prefix := c.String("prefix")
			if os.Getenv("PREFIX") != "" {
				prefix = os.Getenv("PREFIX")
			}

			dir := c.String("dir")
			if os.Getenv("DIR") != "" {
				dir = os.Getenv("DIR")
			}

			var cfg Config
			px, _ := strconv.Atoi(port)
			cfg.Port = int64(px)
			cfg.Prefix = prefix
			cfg.Dir = dir

			Serve(&cfg)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("%s", err.Error())
	}
}
