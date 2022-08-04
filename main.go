package main

import (
	"os"
	"strconv"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/fs"

	"github.com/go-zoox/serve/server"
)

// //go:embed static/*
// var static embed.FS

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
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
	})

	app.Command(func(c *cli.Context) error {
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

		var cfg server.Config
		px, _ := strconv.Atoi(port)
		cfg.Port = int64(px)
		cfg.Prefix = prefix
		cfg.Dir = dir

		// // embed fs
		// cfg.FSMode = "embed"
		// cfg.Dir = "static/"
		// cfg.EmbedFS = &static

		// // proxy
		// cfg.Proxy.Rewrites = map[string]server.ProxyRewrite{
		// 	"^/api/": {
		// 		Target: "http://backend:8080",
		// 		Rewrites: map[string]string{
		// 			"^/api/(.*)": "/$1",
		// 		},
		// 	},
		// 	"^/(.*)$": {
		// 		Target: "http://frontend:8080",
		// 	},
		// }

		server.Serve(&cfg)

		return nil
	})

	app.Run()
}
