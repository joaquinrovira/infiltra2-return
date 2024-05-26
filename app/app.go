package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joaquinrovira/infiltra2-returns/app/services"
	"github.com/samber/do/v2"
)

type App struct {
	Server *http.Server
	Services *do.RootScope
	Context context.Context
}

func NewApp(ctx context.Context) *App {
	services := do.New()
	RegisterAppServices(ctx, services)
	router := AppRouter(services)
	server := NewServer(ctx, router)

	return &App{
		Services: services,
		Server: server,
		Context: ctx,
	}
}

func AppRouter(root *do.RootScope) http.Handler {
	router := chi.NewRouter()
	useDefaultStaticContent(router)
	useRoutes(router, root)
	return router
}

func RegisterAppServices(ctx context.Context,root *do.RootScope) {
	do.Provide(root, func (do.Injector) (context.Context, error) {return ctx, nil})
	do.Provide(root, services.NewRoomsManager)
	do.Provide(root, services.NewRandomWordService)
}

func (a *App) Listen() {
	svr := a.Server
	log.Printf("server listening on %v", svr.Addr)
	if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting server: %s\n", err)
	}
	log.Print("server closed")
}

func (a *App) Shutdown() {
	timeout := 2 * time.Second
	log.Printf("starting graceful shutdown (%v)", timeout)
	ctx, cancel := context.WithTimeout(a.Context, timeout)

	go func() {
		defer cancel()
		a.Server.Shutdown(ctx)
	}()

	<-ctx.Done()

	switch (ctx.Err()) {
	case nil:
	case context.DeadlineExceeded:
	case context.Canceled:
		log.Printf("server shutdown successfully")
	default:
		log.Fatalf("terminated with error: %v", ctx.Err())
	}
}