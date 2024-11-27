package authhandler

import (
	"context"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
)

func (h *AuthHandler) RefreshTokens(ctx context.Context) http.HandlerFunc {
	const op = "handlers.auth.RefreshTokens"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		h.log = with.WithOpAndRequestID(h.log, op, requestID)

		refreshToken := r.Header.Get("Authorization")
		refreshToken = strings.ReplaceAll(refreshToken, "Bearer ", "")
		if refreshToken == "" {
			h.log.Error("missing refresh token")
			ErrorResponse(w, r, 400, "refresh token is required")
			return
		}

		tokens, err := h.authService.RefreshTokens(ctx, refreshToken, r.RemoteAddr, requestID)
		if err != nil {
			h.log.Error("failed to refresh tokens", sl.Err(err))
			ErrorResponse(w, r, 401, "unauthorized")
			return
		}

		SuccessResponse(w, r, 200, tokens)
	}
}
