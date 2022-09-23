package main

import (
	"fmt"

	"github.com/go-zoox/serve/server"
)

func main() {
	cfg := &server.Config{
		Port: 9000,
		//
		Proxy: server.Proxy{
			// 	Rewrites: map[string]server.ProxyRewrite{
			// 		"^/httpbin/": {
			// 			Target: "https://httpbin.zcorky.com",
			// 			Rewrites: map[string]string{
			// 				"^/httpbin/(.*)$": "/$1",
			// 			},
			// 		},
			// 		"^/api/": {
			// 			Target: "http://backend:8080",
			// 			Rewrites: map[string]string{
			// 				"^/api/(.*)$": "/$1",
			// 			},
			// 		},
			// 		"^/(.*)$": {
			// 			Target: "http://frontend:8080",
			// 		},
			// 	},
			// },
			Rewrites: server.ProxyGroupRewrites{
				{
					RegExp: "^/api/",
					Rewrite: server.ProxyRewrite{
						Target: "http://backend:8080",
						Rewrites: server.ProxyRewriters{
							{
								From: "^/api/(.*)$",
								To:   "/$1",
							},
						},
					},
				},
				{
					RegExp: "^/httpbin/",
					Rewrite: server.ProxyRewrite{
						Target: "https://httpbin.zcorky.com",
						Rewrites: server.ProxyRewriters{
							{
								From: "^/httpbin/(.*)$",
								To:   "/$1",
							},
						},
					},
				},
				{
					RegExp: "^/(.*)$",
					Rewrite: server.ProxyRewrite{
						Target: "http://frontend:8080",
					},
				},
			},
		},
	}

	if err := server.Serve(cfg); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
