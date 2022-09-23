package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"

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
			&cli.StringFlag{
				Name:  "basic-auth",
				Usage: "Support basic auth, format: username:password,username2:password2",
			},
			&cli.BoolFlag{
				Name:  "enable-gzip",
				Usage: "Enable gzip compression",
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

		basicAuth := c.String("basic-auth")

		users := map[string]string{}
		if basicAuth != "" {
			for _, u := range strings.Split(basicAuth, ",") {
				userPass := strings.Split(u, ":")
				if len(userPass) != 2 {
					logger.Error("Invalid basic auth user: %s", u)
					continue
				}
				users[userPass[0]] = userPass[1]
			}
		}

		var cfg server.Config
		px, _ := strconv.Atoi(port)
		cfg.Port = int64(px)
		cfg.Prefix = prefix
		cfg.Dir = dir
		cfg.BasicAuth = users
		cfg.EnableGzip = c.Bool("enable-gzip")

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

		return server.Serve(&cfg)
	})

	app.Run()
}
