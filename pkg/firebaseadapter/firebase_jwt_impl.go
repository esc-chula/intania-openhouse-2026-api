package firebaseadapter

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/pkg/jwt"
)

type firebaseJwtImpl struct {
	secretKey []byte
}

func InitFirebaseJwtAdapter(ctx context.Context, secretKey []byte) FirebaseAdapter {
	return &firebaseJwtImpl{
		secretKey: secretKey,
	}
}

func (f *firebaseJwtImpl) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*TokenInfo, error) {
	claims, err := jwt.ParseAuthToken(f.secretKey, idToken)
	if err != nil {
		return nil, err
	}

	uid := getClaimsField(claims, "uid")
	email := getClaimsField(claims, "email")
	display_name := getClaimsField(claims, "name")
	photo_url := getClaimsField(claims, "picture")

	return &TokenInfo{
		UserId:      uid,
		Email:       email,
		DisplayName: display_name,
		PhotoURL:    photo_url,
	}, nil
}
