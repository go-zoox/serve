package main

import "github.com/go-zoox/serve/server"

func main() {
	cfg := &server.Config{
		Port: 9000,
		//
		Proxy: server.Proxy{
			Rewrites: map[string]server.ProxyRewrite{
				"^/api/": {
					Target: "http://backend:8080",
					Rewrites: map[string]string{
						"^/api/(.*)$": "/$1",
					},
				},
				"^/(.*)$": {
					Target: "http://frontend:8080",
				},
			},
		},
	}

	server.Serve(cfg)
}
