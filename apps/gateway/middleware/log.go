package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode codes.Code
}

func (lrw *logResponseWriter) WriteHeader(code codes.Code) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(int(code))
}

func (lrw *logResponseWriter) Unwrap() http.ResponseWriter {
	return lrw.ResponseWriter
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: w}
}

func LoggingMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		lw := newLogResponseWriter(w)
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read request body:", "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		clonedR := r.Clone(r.Context())
		clonedR.Body = io.NopCloser(bytes.NewReader(body))

		next(w, clonedR, pathParams)
		if lw.statusCode != codes.OK {
			slog.Error("Request failed", "method", r.Method, "path", r.URL.Path, "status_code", lw.statusCode, "body", string(body))
		}
	}
}
