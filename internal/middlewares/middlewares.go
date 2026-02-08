package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
)

// TODO:
type Middleware interface{}

type middlewareImpl struct {
	cfg config.Config
	api huma.API
}

func NewMiddleware(cfg config.Config, api huma.API) Middleware {
	return &middlewareImpl{
		cfg: cfg,
		api: api,
	}
}
