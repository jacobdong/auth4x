package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type requestIDKey struct{}

type RequestMiddleware struct{}

func NewRequestMiddleware() *RequestMiddleware {
	return &RequestMiddleware{}
}

func (m *RequestMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		w.Header().Set("X-Request-Id", requestID)

		start := time.Now()
		next(w, r.WithContext(ctx))
		logx.WithContext(ctx).Infof("request_id=%s method=%s path=%s duration=%s", requestID, r.Method, r.URL.Path, time.Since(start))
	}
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(requestIDKey{})
	requestID, ok := value.(string)
	return requestID, ok
}
