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

type TokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLogic {
	return &TokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenLogic) Token(req *types.TokenReq) (resp *types.TokenResp, err error) {
	if req.RefreshToken == "" {
		return nil, errs.New(http.StatusBadRequest, "refresh token is required")
	}

	user, err := l.svcCtx.Store.GetByRefreshToken(req.RefreshToken)
	if err != nil {
		if err == store.ErrInvalidRefreshToken {
			return nil, errs.New(http.StatusUnauthorized, "invalid refresh token")
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
