package firebaseadapter

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"google.golang.org/api/option"
)

type firebaseAuthImpl struct {
	auth *auth.Client
}

func InitFirebaseAuthAdapter(ctx context.Context, cfg config.Config) FirebaseAdapter {
	var opts []option.ClientOption
	if !cfg.App().IsProduction {
		opts = append(opts, option.WithCredentialsFile(cfg.Firebase().ServiceAccountKeyFile))
	}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting firebase auth: %v", err)
	}

	return &firebaseAuthImpl{
		auth: auth,
	}
}

func (f *firebaseAuthImpl) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*TokenInfo, error) {
	token, err := f.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	uid := token.UID
	email := getClaimsField(token.Claims, "email")

	return &TokenInfo{
		UserId: uid,
		Email:  email,
	}, nil
}

func getClaimsField(claims map[string]any, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}
	return ""
}
