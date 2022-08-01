package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	gofs "io/fs"
	"net/http"
	"path"

	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	zd "github.com/go-zoox/zoox/default"
)

// Config is the configuration of the server.
type Config struct {
	Port   int64  `yaml:"port"`
	Prefix string `yaml:"prefix"`
	Dir    string `yaml:"dir"`
	FSMode string `yaml:"fs_mode"` // default: system, optional: system | embed
	//
	EmbedFS *embed.FS
}

// Serve starts the server.
func Serve(cfg *Config) {
	j, _ := json.MarshalIndent(cfg, "", "  ")
	logger.Info(string(j))

	app := zd.Default()

	if cfg.FSMode == "embed" {
		if cfg.EmbedFS == nil {
			panic("fs_mode is embed, but EmbedFS is nil")
		}

		subFS, _ := gofs.Sub(cfg.EmbedFS, "web")

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

	app.Run(fmt.Sprintf(":%d", cfg.Port))
}
