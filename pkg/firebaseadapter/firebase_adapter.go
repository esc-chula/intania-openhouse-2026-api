package firebaseadapter

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInfo struct {
	UserId      string
	Email       string
	DisplayName string
	PhotoURL    string
}

func (t *TokenInfo) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"uid":          t.UserId,
		"email":        t.Email,
		"display_name": t.DisplayName,
		"photo_url":    t.PhotoURL,
	}
}

type FirebaseAdapter interface {
	VerifyIDToken(ctx context.Context, idToken string) (*TokenInfo, error)
}
