package main

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	zd "github.com/go-zoox/zoox/default"
)

type Config struct {
	Port   int64  `yaml:"port"`
	Prefix string `yaml:"prefix"`
	Dir    string `yaml:"dir"`
}

func Serve(cfg *Config) {
	j, _ := json.MarshalIndent(cfg, "", "  ")
	logger.Info(string(j))

	app := zd.Default()

	app.Static(cfg.Prefix, cfg.Dir)

	if indexHTML, err := fs.ReadFileAsString(path.Join(cfg.Dir, "index.html")); err == nil {
		app.Fallback(func(ctx *zoox.Context) {
			ctx.String(200, indexHTML)
		})
	}

	app.Run(fmt.Sprintf(":%d", cfg.Port))
}
