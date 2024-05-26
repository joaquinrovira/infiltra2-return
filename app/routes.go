package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joaquinrovira/infiltra2-returns/app/endpoints"
	"github.com/joaquinrovira/infiltra2-returns/app/routes"
	"github.com/samber/do/v2"
)


func useRoutes(mux *chi.Mux, services *do.RootScope) *chi.Mux {

	
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
