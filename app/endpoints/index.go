package endpoints

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/joaquinrovira/infiltra2-returns/app/components"
)

func Index(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.Home()).ServeHTTP(w,r)
}