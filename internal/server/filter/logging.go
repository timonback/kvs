package filter

import (
	context2 "github.com/timonback/keyvaluestore/internal/server/context"
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logging(Logger *log.Logger, wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(context2.RequestIDKey).(string)

		logRequest := r.URL.Path != context2.HandlerPathInternalId
		if logRequest {
			if !ok {
				requestID = "unknown"
			}
			Logger.Println("-->", requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}

		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)

		if logRequest {
			statusCode := lrw.statusCode
			Logger.Println("<--", requestID, r.Method, r.URL.Path, statusCode, http.StatusText(statusCode))
		}
	})
}
