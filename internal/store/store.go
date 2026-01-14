package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type User struct {
	ID           string
	Username     string
	PasswordHash []byte
	RefreshToken string
}

type UserStore struct {
	mu          sync.RWMutex
	usersByName map[string]*User
	usersByID   map[string]*User
	byRefresh   map[string]*User
}

func NewUserStore() *UserStore {
	return &UserStore{
		usersByName: make(map[string]*User),
		usersByID:   make(map[string]*User),
		byRefresh:   make(map[string]*User),
	}
}

func (s *UserStore) CreateUser(username, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByName[username]; exists {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: hash,
	}

	s.usersByName[username] = user
	s.usersByID[user.ID] = user

	return user, nil
}

func (s *UserStore) Authenticate(username, password string) (*User, error) {
	s.mu.RLock()
	user := s.usersByName[username]
	s.mu.RUnlock()

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserStore) SetRefreshToken(user *User, refreshToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.RefreshToken != "" {
		delete(s.byRefresh, user.RefreshToken)
	}

	user.RefreshToken = refreshToken
	s.byRefresh[refreshToken] = user
}

func (s *UserStore) GetByRefreshToken(token string) (*User, error) {
	s.mu.RLock()
	user := s.byRefresh[token]
	s.mu.RUnlock()

	if user == nil {
		return nil, ErrInvalidRefreshToken
	}

	return user, nil
}

func (s *UserStore) GetByID(id string) (*User, error) {
	s.mu.RLock()
	user := s.usersByID[id]
	s.mu.RUnlock()

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
