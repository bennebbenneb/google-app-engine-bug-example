package main

import (
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify"
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
func getGooglePlaceInfo() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)
		client := urlfetch.Client(ctx)
		resp, err := client.Get(
			"https://maps.googleapis.com/maps/api/place/details/json" +
				"?placeid=[REMOVED]" +
				"&key=[REMOVED]")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		m := minify.New()
		json.Minify(m, w, resp.Body, nil)
	}

	return http.HandlerFunc(fn)
}

func init() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", reactFileServer(fs))
	http.Handle("/services/google-places-info/", getGooglePlaceInfo())
}
