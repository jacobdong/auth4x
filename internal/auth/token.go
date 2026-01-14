package auth

import (
	"errors"
	"time"

	"auth4x/internal/config"
	"auth4x/internal/store"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrInvalidToken = errors.New("invalid token")

func GenerateTokens(cfg config.Config, user *store.User) (accessToken, refreshToken string, expiresIn int64, err error) {
	expireDuration := time.Duration(cfg.Auth.AccessTokenExpireSeconds) * time.Second
	if expireDuration <= 0 {
		expireDuration = time.Hour
	}

	expiresAt := time.Now().Add(expireDuration)
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		return "", "", 0, err
	}

	refreshToken = uuid.NewString()
	expiresIn = int64(expireDuration.Seconds())
	return accessToken, refreshToken, expiresIn, nil
}

func ParseToken(cfg config.Config, tokenStr string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.Auth.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
