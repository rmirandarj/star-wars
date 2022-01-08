package server

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
)

const (
	httpMethodLabelKey  = "method"
	httpCodeLabelKey    = "code"
	httpPathLabelKey    = "path"
	httpKindLabelKey    = "kind"
	levelLabelKey       = "level"
	environmentLabelKey = "environment"
	appNameLabelKey     = "app_name"

	errorLevelLabelValue              = "error"
	successLevelLabelValue            = "success"
	incomingHTTPRequestKindLabelValue = "incoming"
	outgoingHTTPRequestKindLabelValue = "outgoing"
)

var (
	promHTTPRequestsTotalCounter *prometheus.CounterVec
	onceHTTPRequestsTotalCounter sync.Once
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (a App) HTTPServerMetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		method := strings.ToLower(methods[0])

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		httpRequestsTotalIncrement(method, path, strconv.Itoa(rw.statusCode), incomingHTTPRequestKindLabelValue)
	})
}

func httpRequestsTotalIncrement(method, path, code, kind string) {
	getHTTPRequestsTotalCounterInstance().With(prometheus.Labels{
		httpMethodLabelKey: method,
		httpPathLabelKey:   path,
		httpCodeLabelKey:   code,
		httpKindLabelKey:   kind,
	}).Inc()
}

func getHTTPRequestsTotalCounterInstance() *prometheus.CounterVec {
	onceHTTPRequestsTotalCounter.Do(func() {
		promHTTPRequestsTotalCounter = createHTTPRequestsTotalCounter()
	})

	return promHTTPRequestsTotalCounter
}

func createHTTPRequestsTotalCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests by path and HTTP code.",
			ConstLabels: prometheus.Labels{
				environmentLabelKey: viper.GetString("ENVIRONMENT"),
				appNameLabelKey:     viper.GetString("APP_NAME"),
			},
		},
		[]string{httpMethodLabelKey, httpCodeLabelKey, httpPathLabelKey, httpKindLabelKey},
	)
}
