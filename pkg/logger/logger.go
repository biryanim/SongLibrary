package logger

import (
	"go.uber.org/zap"
	"net/http"
)

var Log *zap.Logger = zap.NewNop()

func Initialize() error {
	// Настроим логгер для вывода в консоль в формате JSON
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Log)
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}
