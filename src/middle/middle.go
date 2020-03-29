package middle

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct{}

func (m Middleware) LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("[%s] %q %v", r.Method, r.URL.String(), end.Sub(start))
	}
	return http.HandlerFunc(fn)
}

func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover from panic %+v", err)
				http.Error(w, http.StatusText(
					http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
