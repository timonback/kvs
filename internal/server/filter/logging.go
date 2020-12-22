package filter

import (
	context2 "github.com/timonback/keyvaluestore/internal/server/context"
	"log"
	"net/http"
)

func Logging(Logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(context2.RequestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				Logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
