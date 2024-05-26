package util

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)


func WriteComponent(c templ.Component, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Render(context.TODO(), w)
}