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
	"golang.org/x/crypto/bcrypt"
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
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return nil, errs.BadRequest("username and password are required")
	}

	user, err := l.svcCtx.Store.GetUserByUsername(username)
	if err != nil {
		if err == store.ErrUserNotFound {
			return nil, errs.Unauthorized("invalid credentials")
		}
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errs.Unauthorized("invalid credentials")
	}

	accessToken, refreshToken, expiresIn, err := l.svcCtx.TokenManager.IssueTokens(user, time.Now())
	if err != nil {
		return nil, err
	}

	l.svcCtx.Store.SaveRefreshToken(refreshToken, user.ID, time.Now().Add(l.svcCtx.TokenManager.RefreshExpire), time.Now())

	return &types.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}
