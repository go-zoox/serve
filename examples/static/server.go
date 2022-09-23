package main

import (
	"os"
	"path"

	"github.com/go-zoox/serve/server"
)

func main() {
	pwd, _ := os.Getwd()

	cfg := &server.Config{
		Port:   9000,
		Prefix: "/",
		Dir:    path.Join(pwd, "examples/static/web"),
		FSMode: "system",
		//
		Proxy: server.Proxy{
			Rewrites: server.ProxyGroupRewrites{
				{
					RegExp: "^/api/",
					Rewrite: server.ProxyRewrite{
						Target: "http://backend:8080",
						Rewrites: server.ProxyRewriters{
							{From: "^/api/(.*)$", To: "/$1"},
						},
					},
				},
			},
		},
	}

	server.Serve(cfg)
}
