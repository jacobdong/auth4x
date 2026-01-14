// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"net/http"
	"strings"

	"auth4x/internal/auth"
	"auth4x/internal/errs"
	"auth4x/internal/store"
	"auth4x/internal/svc"
	"auth4x/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req    *http.Request
}

func NewUserMeLogic(ctx context.Context, svcCtx *svc.ServiceContext, req *http.Request) *UserMeLogic {
	return &UserMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req:    req,
	}
}

func (l *UserMeLogic) UserMe() (resp *types.UserResp, err error) {
	authHeader := l.req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errs.New(http.StatusUnauthorized, "missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errs.New(http.StatusUnauthorized, "invalid authorization header")
	}

	claims, err := auth.ParseToken(l.svcCtx.Config, parts[1])
	if err != nil {
		return nil, errs.New(http.StatusUnauthorized, "invalid token")
	}

	user, err := l.svcCtx.Store.GetByID(claims.UserID)
	if err != nil {
		if err == store.ErrUserNotFound {
			return nil, errs.New(http.StatusUnauthorized, "user not found")
		}
		return nil, err
	}

	return &types.UserResp{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
