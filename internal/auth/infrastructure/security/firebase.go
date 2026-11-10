package security

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// FirebaseIdentity thông tin người dùng từ Firebase/Google ID token.
type FirebaseIdentity struct {
	UID         string
	Email       string
	DisplayName string
	AvatarURL   string
}

// TokenVerifier giao diện xác thực Firebase ID token.
type TokenVerifier interface {
	Verify(idToken string) (*FirebaseIdentity, error)
}

// FirebaseTokenVerifier xác thực ID token qua Firebase public keys (giản diện).
// Trong production hãy dùng firebase-admin-sdk-go để verify ký thực tế.
type FirebaseTokenVerifier struct {
	ProjectID string
}

func NewFirebaseTokenVerifier(projectID string) *FirebaseTokenVerifier {
	return &FirebaseTokenVerifier{ProjectID: projectID}
}

// Verify decode payload JWT (không verify ký, chỉ dùng dev/local).
// Production: thay bằng firebase-admin-sdk-go auth.VerifyIDToken.
func (v *FirebaseTokenVerifier) Verify(idToken string) (*FirebaseIdentity, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, err
	}
	id := ""
	if uid, ok := raw["user_id"]; ok {
		id = fmt.Sprintf("%v", uid)
	}
	if id == "" {
		if sub, ok := raw["sub"]; ok {
			id = fmt.Sprintf("%v", sub)
		}
	}
	email := ""
	if e, ok := raw["email"]; ok {
		email = fmt.Sprintf("%v", e)
	}
	name := ""
	if n, ok := raw["name"]; ok {
		name = fmt.Sprintf("%v", n)
	}
	pic := ""
	if p, ok := raw["picture"]; ok {
		pic = fmt.Sprintf("%v", p)
	}
	return &FirebaseIdentity{
		UID:         id,
		Email:       email,
		DisplayName: name,
		AvatarURL:   pic,
	}, nil
}
