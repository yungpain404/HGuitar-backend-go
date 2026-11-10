package common

import (
	"context"
	"net/http"
)

type userIDKey struct{}
type userEmailKey struct{}
type rolesKey struct{}

// WithUserID gán userId vào context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// UserIDFromCtx lấy userId từ context (trống nếu chưa đăng nhập).
func UserIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey{}).(string)
	return v
}

// WithUserEmail gán email vào context.
func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailKey{}, email)
}

// UserEmailFromCtx lấy email từ context.
func UserEmailFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(userEmailKey{}).(string)
	return v
}

// WithRoles gán roles vào context.
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey{}, roles)
}

// RolesFromCtx lấy roles từ context.
func RolesFromCtx(ctx context.Context) []string {
	v, _ := ctx.Value(rolesKey{}).([]string)
	return v
}

// RequireAuth kiểm tra đã đăng nhập chưa; nếu chưa thì ghi lỗi 401 vào w và trả về false.
func RequireAuth(w http.ResponseWriter, r *http.Request) bool {
	if UserIDFromCtx(r.Context()) == "" {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return false
	}
	return true
}
