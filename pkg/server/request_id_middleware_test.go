package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_RequestIdMiddleware(t *testing.T) {
	app := App{}

	testCases := []struct {
		name      string
		requestId string
		handler   func(t *testing.T, writer http.ResponseWriter, request *http.Request)
	}{
		{
			name:      "when request id header with abc23, then it should retrieve the request id from header",
			requestId: "abc23",
			handler: func(t *testing.T, writer http.ResponseWriter, request *http.Request) {
				responseRequestId := writer.Header().Get("x-request-id")
				if responseRequestId != "abc23" {
					t.Fatalf("unexpected response request id, want %s got %s", "abc23", responseRequestId)
				}
			},
		},
		{
			name:      "when request id header not informed, then it should create another one",
			requestId: "",
			handler: func(t *testing.T, writer http.ResponseWriter, request *http.Request) {
				responseRequestId := writer.Header().Get("x-request-id")
				if responseRequestId == "" {
					t.Fatalf("unexpected response request id, want not empty got %s", responseRequestId)
				}
			},
		},
		{
			name:      "when context with request id, then it should retrieve the request id from context",
			requestId: "abc123",
			handler: func(t *testing.T, writer http.ResponseWriter, request *http.Request) {
				requestId := requestIdFromContext(request.Context())
				if requestId != "abc123" {
					t.Fatalf("got unexpected requestId from context, want %s got %s", "abc123", requestId)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/v1/somerequest", nil)
			if tc.requestId != "" {
				r.Header.Set("x-request-id", tc.requestId)
			}
			w := httptest.NewRecorder()

			handler := app.RequestIdMiddleware(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				tc.handler(t, writer, request)
			}))

			handler.ServeHTTP(w, r)

		})

	}
}