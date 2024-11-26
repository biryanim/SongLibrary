package errors

import (
	"encoding/json"
	"github.com/biryanim/SongLibrary/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, errorCode int, errorMessage string) {
	var response ErrorResponse
	response.Code = errorCode
	response.Message = errorMessage
	logger.Log.Debug(response.Message, zap.Any("method", r.Method), zap.Any("url", r.URL.String()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("cannot to marshal", zap.Error(err))
	}
}
