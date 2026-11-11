package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/common"
)

// HomeHandler xử lý các yêu cầu cho trang chủ API.
type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Index trả về thông tin cơ bản về API khi gọi GET /.
func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]any{
		"name":    "HGuitar Backend API",
		"version": "1.0.0",
		"status":  "running",
		"docs":    "/swagger/index.html",
	})
}

// Routes định nghĩa các route cho trang chủ.
func (h *HomeHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.Index)
	return r
}
