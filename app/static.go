package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

const DefaultStaticContentPrefix = "/_/"

func useDefaultStaticContent(mux *chi.Mux) *chi.Mux {
	return useStaticContent(mux, DefaultStaticContentPrefix)
}

func useStaticContent(mux *chi.Mux, pathPrefix string) *chi.Mux {
	workDir, _ := os.Getwd()
	static := http.Dir(filepath.Join(workDir, "static"))
	public := http.Dir(filepath.Join(string(static), "public"))
	
	mux.HandleFunc("GET " + pathPrefix + "*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix(pathPrefix, http.FileServer(public))
		fs.ServeHTTP(w, r)
	})
	return mux
}
