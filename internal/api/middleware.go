package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.URL.Query().Get("request_id")
		if requestID == "" {
			b := make([]byte, 4)
			if _, err := rand.Read(b); err == nil {
				requestID = base64.URLEncoding.EncodeToString(b)[:6]
			} else {
				requestID = time.Now().Format("150405")
			}
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware логирует информацию о запросе
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil {
			b := make([]byte, 4)
			if _, err := rand.Read(b); err == nil {
				requestID = base64.URLEncoding.EncodeToString(b)[:6]
			} else {
				requestID = time.Now().Format("150405")
			}
		}

		next.ServeHTTP(rw, r)

		log.Printf(
			"[%s] %s %s %s %d %v",
			requestID,
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(startTime),
		)
	})
}
