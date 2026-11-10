package common

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

// Recovery middleware ghi log panic và trả về 500.
func Recovery(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

// Logger middleware ghi log request bằng log tiêu chuẩn.
func Logger(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		NoColor: true,
	})(next)
}

// TokenVerifyFunc xác thực JWT và trả về (userID, email, roles, err).
type TokenVerifyFunc func(tokenStr string) (userID string, email string, roles []string, err error)

// AuthMiddleware trích xuất userId từ JWT trong header Authorization: Bearer <token>.
func AuthMiddleware(verify TokenVerifyFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			userID, email, roles, err := verify(tokenStr)
			if err != nil || userID == "" {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			ctx = WithUserID(ctx, userID)
			ctx = WithUserEmail(ctx, email)
			ctx = WithRoles(ctx, roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
