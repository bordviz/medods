package authhandler

import (
	"context"
	"medods/internal/domain/dto"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (h *AuthHandler) CreateUser(ctx context.Context) http.HandlerFunc {
	const op = "handlers.auth.CreateUser"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		h.log = with.WithOpAndRequestID(h.log, op, requestID)

		var user dto.User
		if err := render.Decode(r, &user); err != nil {
			h.log.Error("failed to decode model", sl.Err(err))
			ErrorResponse(w, r, 400, "Bad request")
			return
		}
		if err := user.Validate(); err != nil {
			h.log.Error("validate model error", sl.Err(err))
			ErrorResponse(w, r, 422, err.Error())
			return
		}

		userID, err := h.authService.CreateUser(ctx, &user, requestID)
		if err != nil {
			h.log.Error("failed to create user", sl.Err(err))
			ErrorResponse(w, r, 400, err.Error())
			return
		}

		SuccessResponse(w, r, 201, map[string]string{
			"detail": "new user was successfully created",
			"id":     userID.String(),
		})
	}
}
