package services

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
)

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
	http.Handle("/services/google-places-info/", getGooglePlaceInfo())
}
