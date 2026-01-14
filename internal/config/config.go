// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth AuthConfig `json:"Auth"`
}

type AuthConfig struct {
	JWTSecret                string `json:"JWTSecret"`
	AccessTokenExpireSeconds int64  `json:"AccessTokenExpireSeconds"`
}
