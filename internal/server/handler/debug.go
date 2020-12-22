package handler

import (
	"fmt"
	"net/http"
)

func Debug() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-INFO", r.RequestURI)
		fmt.Fprintln(w, r.RequestURI)
	})
}
