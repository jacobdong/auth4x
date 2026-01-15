// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"strings"
	"time"

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
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, errs.BadRequest("refresh token is required")
	}

	claims, err := l.svcCtx.TokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errs.Unauthorized("invalid refresh token")
	}

	storedToken, err := l.svcCtx.Store.GetRefreshToken(refreshToken)
	if err != nil {
		if err == store.ErrRefreshTokenAbsent {
			return nil, errs.Unauthorized("refresh token not found")
		}
		return nil, err
	}

	if storedToken.RevokedAt != nil {
		return nil, errs.Unauthorized("refresh token revoked")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, errs.Unauthorized("refresh token expired")
	}

	user, err := l.svcCtx.Store.GetUserByID(claims.Subject)
	if err != nil {
		return nil, errs.Unauthorized("user not found")
	}

	l.svcCtx.Store.RevokeRefreshToken(refreshToken, time.Now())

	accessToken, newRefreshToken, expiresIn, err := l.svcCtx.TokenManager.IssueTokens(user, time.Now())
	if err != nil {
		return nil, err
	}

	l.svcCtx.Store.SaveRefreshToken(newRefreshToken, user.ID, time.Now().Add(l.svcCtx.TokenManager.RefreshExpire), time.Now())

	return &types.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}
