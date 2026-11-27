package httpserver

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

const (
	immutableCacheControl = "public, max-age=31536000, immutable"
	htmlCacheControl      = "no-store, no-cache, must-revalidate"
	manifestCacheControl  = "no-cache"
)

// NewStaticUIHandler returns an http.Handler serving a static-exported SPA from
// fsys. It resolves exact regular files, <path>/index.html, <path>.html, and
// falls back to root index.html for non-asset client-side routes. Immutable
// Cache-Control headers are applied to /assets/* (Vite) and /_next/static/*
// hashed assets, and no-store headers are applied to HTML documents.
func NewStaticUIHandler(fsys fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if fsys == nil {
			http.NotFound(w, r)
			return
		}

		if hasUnsafePath(r.URL.Path) || (r.URL.RawPath != "" && hasUnsafePath(r.URL.RawPath)) {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		cleaned := path.Clean("/" + r.URL.Path)
		rel := strings.TrimPrefix(cleaned, "/")
		if rel == "." {
			rel = ""
		}

		servedName, file, info, ok := resolveStaticFile(fsys, rel)
		if !ok {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		applyStaticCacheHeaders(w.Header(), servedName)
		serveFSFile(w, r, servedName, file, info)
	})
}

func hasUnsafePath(raw string) bool {
	if strings.ContainsRune(raw, '\x00') || strings.ContainsRune(raw, '\\') {
		return true
	}
	for _, seg := range strings.Split(raw, "/") {
		if seg == ".." {
			return true
		}
	}
	return false
}

func resolveStaticFile(fsys fs.FS, rel string) (string, fs.File, fs.FileInfo, bool) {
	if rel != "" {
		if f, info, ok := openRegularFile(fsys, rel); ok {
			return rel, f, info, true
		}
	}

	dirIndex := path.Join(rel, "index.html")
	if f, info, ok := openRegularFile(fsys, dirIndex); ok {
		return dirIndex, f, info, true
	}

	if rel != "" {
		htmlCandidate := rel + ".html"
		if f, info, ok := openRegularFile(fsys, htmlCandidate); ok {
			return htmlCandidate, f, info, true
		}
	}

	if !isMissingAssetPath(rel) {
		if f, info, ok := openRegularFile(fsys, "index.html"); ok {
			return "index.html", f, info, true
		}
	}

	return "", nil, nil, false
}

func openRegularFile(fsys fs.FS, name string) (fs.File, fs.FileInfo, bool) {
	if name == "" || !fs.ValidPath(name) {
		return nil, nil, false
	}
	f, err := fsys.Open(name)
	if err != nil {
		return nil, nil, false
	}
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		_ = f.Close()
		return nil, nil, false
	}
	return f, info, true
}

func isMissingAssetPath(rel string) bool {
	if rel == "" {
		return false
	}
	if rel == "assets" || strings.HasPrefix(rel, "assets/") || rel == "_next" || strings.HasPrefix(rel, "_next/") {
		return true
	}
	return path.Ext(path.Base(rel)) != ""
}

func applyStaticCacheHeaders(h http.Header, servedName string) {
	lower := strings.ToLower(servedName)
	switch {
	case strings.HasPrefix(servedName, "assets/"), strings.HasPrefix(servedName, "_next/static/"):
		h.Set("Cache-Control", immutableCacheControl)
	case lower == "sw.js":
		h.Set("Cache-Control", htmlCacheControl)
		h.Set("Service-Worker-Allowed", "/")
		h.Set("Content-Type", "text/javascript; charset=utf-8")
	case strings.HasSuffix(lower, ".webmanifest"):
		h.Set("Cache-Control", manifestCacheControl)
		h.Set("Content-Type", "application/manifest+json")
	case strings.HasSuffix(lower, ".html"):
		h.Set("Cache-Control", htmlCacheControl)
	}
}

func serveFSFile(w http.ResponseWriter, r *http.Request, servedName string, file fs.File, info fs.FileInfo) {
	name := path.Base(servedName)
	if rs, ok := file.(io.ReadSeeker); ok {
		http.ServeContent(w, r, name, info.ModTime(), rs)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, name, info.ModTime(), bytes.NewReader(data))
}
