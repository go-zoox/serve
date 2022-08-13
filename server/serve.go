package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	gofs "io/fs"
	"net/http"
	"path"

	"github.com/go-zoox/debug"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	zd "github.com/go-zoox/zoox/default"
	"github.com/go-zoox/zoox/middleware"
)

// Config is the configuration of the server.
type Config struct {
	Port   int64  `yaml:"port"`
	Prefix string `yaml:"prefix"`
	Dir    string `yaml:"dir"`
	FSMode string `yaml:"fs_mode"` // default: system, optional: system | embed
	//
	EmbedFS *embed.FS
	//
	Proxy Proxy `yaml:"proxy"`
	// Basic Auth Users Map
	BasicAuth map[string]string `yaml:"basic_auth"`
	//
	EnableGzip bool
}

// Proxy is the proxy configuration.
type Proxy struct {
	// Rewrites Example:
	//	{
	//		"^/api": {
	//			target: "http://backend:8080",
	//			rewrites: {
	//				"^/api/(.*)$": "/$1",
	//			},
	//		},
	//	}
	Rewrites map[string]ProxyRewrite `yaml:"rewrites"`
}

// ProxyRewrite is the proxy rewrite target configuration.
type ProxyRewrite struct {
	Target   string            `yaml:"target"`
	Rewrites map[string]string `yaml:"rewrites"`
}

// Serve starts the server.
func Serve(cfg *Config) error {
	if debug.IsDebugMode() {
		j, _ := json.MarshalIndent(cfg, "", "  ")
		logger.Info(string(j))
	}

	app := zd.Default()

	if cfg.EnableGzip {
		app.Use(middleware.Gzip())
	}

	if len(cfg.BasicAuth) > 0 {
		app.Use(middleware.BasicAuth("go-zoox/serve", cfg.BasicAuth))
	}

	if cfg.Proxy.Rewrites != nil {
		rewrites := make(map[string]middleware.ProxyRewrite)
		for k, v := range cfg.Proxy.Rewrites {
			rewrites[k] = middleware.ProxyRewrite{
				Target:   v.Target,
				Rewrites: v.Rewrites,
			}
		}

		app.Use(middleware.Proxy(&middleware.ProxyConfig{
			Rewrites: rewrites,
		}))
	}

	if cfg.FSMode == "embed" {
		if cfg.EmbedFS == nil {
			panic("fs_mode is embed, but EmbedFS is nil")
		}

		if cfg.Dir == "" {
			panic("dir is required, but empty")
		}

		if cfg.Dir[0] == '/' {
			panic("dir must be relative path in embed mode")
		}

		subFS, _ := gofs.Sub(cfg.EmbedFS, cfg.Dir)

		app.StaticFS(cfg.Prefix, http.FS(subFS))

		if f, err := subFS.Open("index.html"); err == nil {
			if bytes, err := io.ReadAll(f); err == nil {
				app.Fallback(func(ctx *zoox.Context) {
					ctx.Status(200)
					ctx.Write(bytes)
				})
			}
		}
	} else {
		app.Static(cfg.Prefix, cfg.Dir)

		if indexHTML, err := fs.ReadFileAsString(path.Join(cfg.Dir, "index.html")); err == nil {
			app.Fallback(func(ctx *zoox.Context) {
				ctx.String(200, indexHTML)
			})
		}
	}

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
