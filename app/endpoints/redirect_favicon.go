package endpoints

import (
	"net/http"
)

func RedirectFavicon(route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route, http.StatusMovedPermanently)
	}
}