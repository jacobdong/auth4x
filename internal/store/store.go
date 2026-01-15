package store

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrRefreshTokenAbsent = errors.New("refresh token not found")
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type Store struct {
	mu            sync.RWMutex
	usersByID     map[string]User
	usersByName   map[string]User
	refreshTokens map[string]RefreshToken
}

func NewStore() *Store {
	return &Store{
		usersByID:     make(map[string]User),
		usersByName:   make(map[string]User),
		refreshTokens: make(map[string]RefreshToken),
	}
}

func (s *Store) CreateUser(id, username, passwordHash string, createdAt time.Time) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByName[username]; exists {
		return User{}, ErrUsernameExists
	}

	user := User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
	s.usersByID[id] = user
	s.usersByName[username] = user
	return user, nil
}

func (s *Store) GetUserByUsername(username string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByName[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *Store) GetUserByID(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *Store) SaveRefreshToken(token string, userID string, expiresAt time.Time, createdAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshTokens[token] = RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}
}

func (s *Store) GetRefreshToken(token string) (RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	refreshToken, ok := s.refreshTokens[token]
	if !ok {
		return RefreshToken{}, ErrRefreshTokenAbsent
	}
	return refreshToken, nil
}

func (s *Store) RevokeRefreshToken(token string, revokedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	refreshToken, ok := s.refreshTokens[token]
	if !ok {
		return
	}
	refreshToken.RevokedAt = &revokedAt
	s.refreshTokens[token] = refreshToken
}
