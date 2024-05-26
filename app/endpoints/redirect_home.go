package endpoints

import (
	"net/http"
)

func RedirectHome(route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, route, http.StatusSeeOther)
	}
}