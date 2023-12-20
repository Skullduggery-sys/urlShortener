package httpUtils

import (
	"encoding/json"
	"net/http"
	"urlShortener/utils/e"
)

func RenderJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	const fn = "http.httpUtils.RenderJSON"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "rendering error", http.StatusInternalServerError)
		return e.WrapError(fn, err)
	}
	return nil
}
