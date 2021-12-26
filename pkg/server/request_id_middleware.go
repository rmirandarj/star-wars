package server

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	xRequestIdHeader = "x-request-id"
)
func (a App) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(xRequestIdHeader)
		if requestId == "" {
			requestId = uuid.NewString()
		}
		w.Header().Set(xRequestIdHeader, requestId)
		ctx := withRequestId(r.Context(), requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, xRequestIdHeader, requestId)
}

func requestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(xRequestIdHeader).(string)
	if !ok {
		return ""
	}
	return requestId
}
