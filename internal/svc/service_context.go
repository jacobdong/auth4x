// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"auth4x/internal/auth"
	"auth4x/internal/config"
	"auth4x/internal/store"
)

type ServiceContext struct {
	Config       config.Config
	Store        *store.Store
	TokenManager *auth.TokenManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		Store:        store.NewStore(),
		TokenManager: auth.NewTokenManager(c.Auth.AccessSecret, c.Auth.AccessExpire, c.Auth.RefreshSecret, c.Auth.RefreshExpire),
	}
}
