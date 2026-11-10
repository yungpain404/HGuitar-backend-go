package common

import (
	"encoding/json"
	"net/http"
)

// WriteJSON ghi JSON response.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError ghi lỗi AppError dạng JSON.
func WriteError(w http.ResponseWriter, err error) {
	appErr := FromError(err)
	WriteJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
}

// DecodeJSON giải mã body JSON vào v.
func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
