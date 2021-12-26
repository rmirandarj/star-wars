package server

import (
	"context"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func configureLog(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warn(err)
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func loggerFromRequest(req *http.Request) *log.Entry {
	return loggerFromContext(req.Context())
}

func loggerFromContext(ctx context.Context) *log.Entry {
	return log.WithField("x-request-id", requestIdFromContext(ctx))
}
