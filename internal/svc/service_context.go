// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"auth4x/internal/config"
	"auth4x/internal/store"
)

type ServiceContext struct {
	Config config.Config
	Store  *store.UserStore
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Store:  store.NewUserStore(),
	}
}
