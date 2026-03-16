package server

import (
	"context"
	"errors"
	"log"

	"github.com/danielgtaylor/huma/v2"
)

type handlerErrorContextKey struct{}

type handlerErrorState struct {
	lastError error
}

// ErrorCaptureTransformer is a Huma transformer that captures errors from handler responses.
func ErrorCaptureTransformer(ctx huma.Context, status string, v any) (any, error) {
	if em, ok := v.(*huma.ErrorModel); ok {
		if em.Status >= 500 {
			// Create a string of all error messages to log
			errMsg := em.Detail
			if len(em.Errors) > 0 {
				errMsg += ": "
				for i, e := range em.Errors {
					if i > 0 {
						errMsg += ", "
					}
					errMsg += e.Message
				}
			}
			recordHandlerError(ctx.Context(), errors.New(errMsg))

			// Sanitize the error model to avoid leaking sensitive internal details
			em.Detail = "internal server error"
			em.Errors = nil
		} else {
			recordHandlerError(ctx.Context(), em)
		}
	} else if err, isError := v.(error); isError {
		recordHandlerError(ctx.Context(), err)
	}
	return v, nil
}

// ErrorRecorderMiddleware adds error recording capability to the request context.
// This middleware must be applied before any handlers that need error recording.
func ErrorRecorderMiddleware(ctx huma.Context, next func(huma.Context)) {
	errorState := &handlerErrorState{}
	newCtx := huma.WithValue(ctx, handlerErrorContextKey{}, errorState)
	next(newCtx)
}

// ErrorLoggerMiddleware logs the captured errors.
// This middleware should be applied after ErrorRecorderMiddleware.
func ErrorLoggerMiddleware(ctx huma.Context, next func(huma.Context)) {
	defer func() {
		if err := getHandlerError(ctx.Context()); err != nil {
			log.Println("REQUEST ERROR:", err)
		}
	}()
	next(ctx)
}

// getHandlerError retrieves the last error recorded by the ErrorCaptureTransformer.
// Returns nil if no error was recorded or if ErrorRecorderMiddleware wasn't applied.
func getHandlerError(ctx context.Context) error {
	state := ctx.Value(handlerErrorContextKey{})
	if state == nil {
		return nil
	}

	errorState, ok := state.(*handlerErrorState)
	if !ok {
		return nil
	}

	return errorState.lastError
}

// recordHandlerError stores an error in the request context.
// This is used internally by ErrorCaptureTransformer.
func recordHandlerError(ctx context.Context, err error) {
	state := ctx.Value(handlerErrorContextKey{})
	if state == nil {
		return
	}

	errorState, ok := state.(*handlerErrorState)
	if !ok {
		return
	}

	errorState.lastError = err
}
