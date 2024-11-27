package jwt

import (
	"encoding/base64"
	"fmt"
	"medods/internal/config"
	"medods/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type JWTAuth struct {
	cfg *config.Auth
}

func NewJWTAuth(cfg *config.Auth) *JWTAuth {
	return &JWTAuth{cfg: cfg}
}

func (j *JWTAuth) createAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(j.cfg.AccessTokenLifetime).Unix()

	tokenString, err := token.SignedString([]byte(j.cfg.AccessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTAuth) createRefreshToken(userID uuid.UUID, ipAddress string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["ip_address"] = ipAddress
	claims["exp"] = time.Now().Add(j.cfg.RefreshTokenLifetime).Unix()

	tokenString, err := token.SignedString([]byte(j.cfg.RefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTAuth) CreateTokenPair(userID uuid.UUID, ipAddress string) (*models.TokenPair, error) {
	accessToken, err := j.createAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.createRefreshToken(userID, ipAddress)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWTAuth) DecodeAccessToken(token string) (*models.AccessToken, error) {
	var model models.AccessToken

	jwtToken, err := jwt.ParseWithClaims(token, &model, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.AccessSecret), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return &model, nil
}

func (j *JWTAuth) DecodeRefreshToken(token string) (*models.RefreshToken, error) {
	var model models.RefreshToken

	jwtToken, err := jwt.ParseWithClaims(token, &model, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.RefreshSecret), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return &model, nil
}

func (j JWTAuth) HashRefreshToken(token string) string {
	salt := []byte("hdaBUYgasdjksgTUYs")
	hash := argon2.IDKey([]byte(token), salt, 1, 64*1024, 4, 32)
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	return encodedHash
}
