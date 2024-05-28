package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/joaquinrovira/infiltra2-returns/app/components"
	"github.com/joaquinrovira/infiltra2-returns/app/endpoints"
	"github.com/joaquinrovira/infiltra2-returns/app/routes"
	"github.com/samber/do/v2"
)


func useRoutes(mux *chi.Mux, services *do.RootScope) *chi.Mux {

	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) { templ.Handler(components.Test()).ServeHTTP(w,r) })
	
	mux.HandleFunc("GET " + routes.Home(), endpoints.Index)
	mux.HandleFunc("GET " + routes.Lobby(), endpoints.Lobby())
	
	mux.HandleFunc("GET " + routes.RoomTemplate(), endpoints.Room(services))
	mux.HandleFunc("GET " + routes.RoomSSETemplate(), endpoints.RoomEvents(services))
	
	mux.HandleFunc(routes.ReadyTemplate(), endpoints.Ready(services))
	
	// Misc stuff
	mux.HandleFunc("GET " + routes.CatchAll(), endpoints.RedirectHome(routes.Home()))
	mux.HandleFunc("GET /favicon.ico", endpoints.RedirectFavicon(DefaultStaticContentPrefix + "favicon.ico"))
	
	return mux
}
