//go:build !embedui

package main

import (
	"io/fs"
	"net/http"
)

const devUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>omeg4mon-router — Dev Mode</title>
</head>
<body>
  <h1>omeg4mon-router (dev mode)</h1>
  <p>The frontend bundle is not embedded in this build.</p>
  <p>Run the Vite dev server and visit <a href="http://127.0.0.1:20127">http://127.0.0.1:20127</a>, or build the single binary with <code>make build</code>.</p>
</body>
</html>
`

func loadUI() (fs.FS, http.Handler, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(devUIHTML))
		}
	})
	return nil, handler, nil
}
