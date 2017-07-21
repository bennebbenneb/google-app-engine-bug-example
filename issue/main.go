package main

import (
	"net/http"
	"strings"
)

func reactFileServer(fs http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		staticPaths := [...]string{
			"/static/",
			"/img/",
		}

		var isStaticPath bool = false
		for _, path := range staticPaths {
			if strings.HasPrefix(req.URL.Path, path) {
				isStaticPath = true
				break
			}
		}

		if isStaticPath {
			fs.ServeHTTP(w, req)
		} else {
			fsHandler := http.StripPrefix(req.URL.Path, fs)
			fsHandler.ServeHTTP(w, req)
		}
	}
	return http.HandlerFunc(fn)
}

func init() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", reactFileServer(fs))
}
