package auth

import (
	"errors"
	"time"

	"auth4x/internal/store"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	Username  string `json:"username"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	AccessSecret  []byte
	AccessExpire  time.Duration
	RefreshSecret []byte
	RefreshExpire time.Duration
}

func NewTokenManager(accessSecret string, accessExpireSeconds int64, refreshSecret string, refreshExpireSeconds int64) *TokenManager {
	return &TokenManager{
		AccessSecret:  []byte(accessSecret),
		AccessExpire:  time.Duration(accessExpireSeconds) * time.Second,
		RefreshSecret: []byte(refreshSecret),
		RefreshExpire: time.Duration(refreshExpireSeconds) * time.Second,
	}
}

func (m *TokenManager) IssueTokens(user store.User, now time.Time) (string, string, int64, error) {
	accessExp := now.Add(m.AccessExpire)
	refreshExp := now.Add(m.RefreshExpire)

	accessToken, err := m.issueToken(user, TokenTypeAccess, accessExp, m.AccessSecret)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := m.issueToken(user, TokenTypeRefresh, refreshExp, m.RefreshSecret)
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, int64(m.AccessExpire.Seconds()), nil
}

func (m *TokenManager) issueToken(user store.User, tokenType string, expiresAt time.Time, secret []byte) (string, error) {
	claims := Claims{
		Username:  user.Username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (m *TokenManager) ParseAccessToken(tokenString string) (*Claims, error) {
	return m.parseToken(tokenString, TokenTypeAccess, m.AccessSecret)
}

func (m *TokenManager) ParseRefreshToken(tokenString string) (*Claims, error) {
	return m.parseToken(tokenString, TokenTypeRefresh, m.RefreshSecret)
}

func (m *TokenManager) parseToken(tokenString string, expectedType string, secret []byte) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
