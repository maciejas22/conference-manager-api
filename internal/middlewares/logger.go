package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"Method: %s, URI: %s, User Agent: %s, Remote Addr: %s, Time: %s",
			r.Method,
			r.RequestURI,
			r.UserAgent(),
			r.RemoteAddr,
			start.Format(time.RFC1123),
		)
		next.ServeHTTP(w, r)
	})
}
