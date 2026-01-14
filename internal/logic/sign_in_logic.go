// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"net/http"

	"auth4x/internal/auth"
	"auth4x/internal/errs"
	"auth4x/internal/store"
	"auth4x/internal/svc"
	"auth4x/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInLogic {
	return &SignInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInLogic) SignIn(req *types.SignInReq) (resp *types.TokenResp, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errs.New(http.StatusBadRequest, "username and password are required")
	}

	user, err := l.svcCtx.Store.Authenticate(req.Username, req.Password)
	if err != nil {
		if err == store.ErrInvalidCredentials {
			return nil, errs.New(http.StatusUnauthorized, "invalid credentials")
		}
		return nil, err
	}

	accessToken, refreshToken, expiresIn, err := auth.GenerateTokens(l.svcCtx.Config, user)
	if err != nil {
		return nil, err
	}

	l.svcCtx.Store.SetRefreshToken(user, refreshToken)

	return &types.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}
