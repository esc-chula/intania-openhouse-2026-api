package firebaseadapter

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInfo struct {
	UserId string
	Email  string
}

func (t *TokenInfo) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"uid":   t.UserId,
		"email": t.Email,
	}
}

type FirebaseAdapter interface {
	VerifyIDToken(ctx context.Context, idToken string) (*TokenInfo, error)
}
