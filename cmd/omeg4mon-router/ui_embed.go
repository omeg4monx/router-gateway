//go:build embedui

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed all:ui_dist
var embeddedUI embed.FS

func loadUI() (fs.FS, http.Handler, error) {
	sub, err := fs.Sub(embeddedUI, "ui_dist")
	if err != nil {
		return nil, nil, fmt.Errorf("open embedded ui_dist: %w", err)
	}
	return sub, nil, nil
}
