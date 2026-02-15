package middlewares

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/firebaseadapter"
)
// Middleware interface
type Middleware interface {
	WithAuthContext(ctx huma.Context, next func(huma.Context))
}

type middlewareImpl struct {
	cfg             config.Config
	api             huma.API
	firebaseAdapter firebaseadapter.FirebaseAdapter
}

func NewMiddleware(cfg config.Config, api huma.API, firebaseAdapter firebaseadapter.FirebaseAdapter) Middleware {
	return &middlewareImpl{
		cfg:             cfg,
		api:             api,
		firebaseAdapter: firebaseAdapter,
	}
}

func (m *middlewareImpl) WithAuthContext(ctx huma.Context, next func(huma.Context)) {
	header := ctx.Header("Authorization")
	if header == "" {
		huma.WriteErr(m.api, ctx, http.StatusBadRequest, "Authorization header not found")
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Invalid Authorization header format")
		return
	}
	idToken := parts[1]

	tokenInfo, err := m.firebaseAdapter.VerifyIDToken(ctx.Context(), idToken)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	ctx = huma.WithValue(ctx, "uid", tokenInfo.UserId)
	ctx = huma.WithValue(ctx, "email", tokenInfo.Email)

	next(ctx)
}
