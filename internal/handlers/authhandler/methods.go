package authhandler

import (
	"context"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/handlers/response"
	"medods/internal/lib/customerror"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// @Summary Create a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param		body 	body	dto.User		true	"Request body"
// @Success		201			{object}	response.Response		"success response"
// @Failure		400,422,500,401			{object}	response.Response	"failure response"
// @Router /auth/create_user [post]
func (h *AuthHandler) CreateUser() http.HandlerFunc {
	const op = "handlers.auth.CreateUser"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		log := with.WithOpAndRequestID(h.log, op, requestID)

		var user dto.User
		if err := render.Decode(r, &user); err != nil {
			log.Error("failed to decode model", sl.Err(err))
			response.ErrorResponse(w, r, response.ErrBadRequest)
			return
		}
		if cerr := user.Validate(); cerr != nil {
			log.Error("validate model error", sl.Err(cerr))
			response.ErrorResponse(w, r, cerr)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userID, cerr := h.authService.CreateUser(ctx, &user, requestID)
		if cerr != nil {
			log.Error("failed to create user", sl.Err(cerr))
			response.ErrorResponse(w, r, cerr)
			return
		}

		response.SuccessResponse(w, r, 201, models.CreateUser{
			Detail: "new user was successfully created",
			ID:     userID.String(),
		})
	}
}

// @Summary Create token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param 		userID 		path 		string 		true 		"User ID"
// @Success		200			{object}	models.TokenPair		"Token pair"
// @Failure		400,422,500,401			{object}	response.Response	"failure response"
// @Router /auth/get_tokens/{userID} [get]
func (h *AuthHandler) CreateTokenPair() http.HandlerFunc {
	const op = "handlers.auth.GetTokens"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		log := with.WithOpAndRequestID(h.log, op, requestID)

		idString := chi.URLParam(r, "userID")
		userID, err := uuid.Parse(idString)
		if err != nil {
			log.Error("failed to parse user id", sl.Err(err))
			response.ErrorResponse(w, r, response.ErrBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		tokens, cerr := h.authService.CreateTokenPair(ctx, userID, ipAddress, requestID)
		if cerr != nil {
			log.Error("failed to get tokens", sl.Err(cerr))
			response.ErrorResponse(w, r, cerr)
			return
		}

		response.SuccessResponse(w, r, 200, tokens)
	}
}

// @Summary Refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param		Authorization header string true "Refresh token"
// @Success		200			{object}	models.TokenPair		"Token pair"
// @Failure		400,422,500,401			{object}	response.Response	"failure response"
// @Router /auth/refresh_tokens [get]
func (h *AuthHandler) RefreshTokens() http.HandlerFunc {
	const op = "handlers.auth.RefreshTokens"

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		log := with.WithOpAndRequestID(h.log, op, requestID)

		refreshToken := r.Header.Get("Authorization")
		refreshToken = strings.ReplaceAll(refreshToken, "Bearer ", "")
		if refreshToken == "" {
			log.Error("missing refresh token")
			response.ErrorResponse(w, r, customerror.NewCustomError("refresh token is required", 400))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		tokens, cerr := h.authService.RefreshTokens(ctx, refreshToken, ipAddress, requestID)
		if cerr != nil {
			log.Error("failed to refresh tokens", sl.Err(cerr))
			response.ErrorResponse(w, r, customerror.NewCustomError("unauthorized", 401))
			return
		}

		response.SuccessResponse(w, r, 200, tokens)
	}
}
