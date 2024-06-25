package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	gofs "io/fs"
	"net/http"
	"path"
	"regexp"

	"github.com/go-zoox/debug"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
	defaults "github.com/go-zoox/zoox/defaults"
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
	//
	Middlewares []zoox.Middleware
	//
	Api              string
	ApiPath          string
	IsApiPathRewrite bool
}

// Proxy is the proxy configuration.
type Proxy struct {
	// Rewrites Example:
	//	[
	//		{
	//			regexp: "^/api",
	//			rewrite: []{
	//				target: "http://backend:8080",
	//				rewrites: []{
	//					from: "^/api/(.*)$",
	//					to: "/$1",
	//				}
	//			},
	//		},
	//	]
	Rewrites ProxyGroupRewrites `yaml:"rewrites"`

	// Dynamic Example:
	Dynamic func(ctx *zoox.Context) ProxyGroupRewrites
}

// ProxyGroupRewrites is the proxy group configuration.
type ProxyGroupRewrites = middleware.ProxyGroupRewrites

// ProxyRewrite is the proxy rewrite configuration.
type ProxyRewrite = middleware.ProxyRewrite

// ProxyRewriters is the proxy rewriters.
type ProxyRewriters = rewriter.Rewriters

// Serve starts the server.
func Serve(cfg *Config) error {
	if debug.IsDebugMode() {
		j, _ := json.MarshalIndent(cfg, "", "  ")
		logger.Info(string(j))
	}

	app := defaults.Default()

	// app.Config.HTTPSPort = 8443
	// app.Config.TLSCertFile = "/opt/data/plugins/local-https/domain/localhost/server.crt"
	// app.Config.TLSKeyFile = "/opt/data/plugins/local-https/domain/localhost/server.key"

	if cfg.EnableGzip {
		app.Use(middleware.Gzip())
	}

	if len(cfg.BasicAuth) > 0 {
		app.Use(middleware.BasicAuth("go-zoox/serve", cfg.BasicAuth))
	}

	if cfg.Middlewares != nil {
		app.Use(cfg.Middlewares...)
	}

	if cfg.Proxy.Dynamic != nil {
		app.Use(func(ctx *zoox.Context) {
			rewrites := cfg.Proxy.Dynamic(ctx)

			for _, group := range rewrites {
				if matched, err := regexp.MatchString(group.RegExp, ctx.Path); err == nil && matched {
					// @BUG: this is not working
					p := proxy.NewSingleHost(group.Rewrite.Target, &proxy.SingleHostConfig{
						Rewrites: group.Rewrite.Rewrites,
						// ChangeOrigin: true,
					})

					p.ServeHTTP(ctx.Writer, ctx.Request)
					return
				}
			}

			ctx.Next()
		})
	}

	if cfg.Proxy.Rewrites != nil {
		app.Use(middleware.ProxyGroups(&middleware.ProxyGroupsConfig{
			Rewrites: cfg.Proxy.Rewrites,
		}))
	}

	if cfg.Api != "" {
		to := fmt.Sprintf("%s/$1", cfg.ApiPath)
		if cfg.IsApiPathRewrite {
			to = "/$1"
		}

		app.Proxy(cfg.ApiPath, cfg.Api, func(cfgX *zoox.ProxyConfig) {
			cfgX.Rewrites = ProxyRewriters{
				{
					From: fmt.Sprintf("%s/(.*)", cfg.ApiPath),
					To:   to,
				},
			}

			cfgX.OnRequest = func(req *http.Request) error {
				return nil
			}
		})
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
