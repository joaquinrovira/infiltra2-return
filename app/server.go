package main

import (
	"context"
	"net"
	"net/http"
	"os"
)

var LISTEN_ADDR = ":8080"

func init() {
	envListenAddr := os.Getenv("LISTEN_ADDR")
	if envListenAddr != "" {
		LISTEN_ADDR = envListenAddr
	}
}

func NewServer(ctx context.Context, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr: LISTEN_ADDR,
		Handler: handler,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}
	return srv
}
