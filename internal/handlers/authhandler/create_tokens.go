package authhandler

import (
	"context"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

func (h *AuthHandler) CreateTokenPair(ctx context.Context) http.HandlerFunc {
	const op = "handlers.auth.GetTokens"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		h.log = with.WithOpAndRequestID(h.log, op, requestID)

		idString := r.URL.Query().Get("userID")
		userID, err := uuid.Parse(idString)
		if err != nil {
			h.log.Error("failed to parse user id", sl.Err(err))
			ErrorResponse(w, r, 400, "Bad request")
			return
		}

		tokens, err := h.authService.CreateTokenPair(ctx, userID, r.RemoteAddr, requestID)
		if err != nil {
			h.log.Error("failed to get tokens", sl.Err(err))
			ErrorResponse(w, r, 400, err.Error())
			return
		}

		SuccessResponse(w, r, 200, tokens)
	}
}
