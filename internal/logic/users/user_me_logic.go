// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package users

import (
	"context"
	"strings"

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
}

func NewUserMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMeLogic {
	return &UserMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserMeLogic) UserMe(req *types.UserMeReq) (resp *types.UserResp, err error) {
	authHeader := strings.TrimSpace(req.Authorization)
	if authHeader == "" {
		return nil, errs.Unauthorized("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errs.Unauthorized("invalid authorization header")
	}

	claims, err := l.svcCtx.TokenManager.ParseAccessToken(parts[1])
	if err != nil {
		return nil, errs.Unauthorized("invalid access token")
	}

	user, err := l.svcCtx.Store.GetUserByID(claims.Subject)
	if err != nil {
		if err == store.ErrUserNotFound {
			return nil, errs.NotFound("user not found")
		}
		return nil, err
	}

	return &types.UserResp{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
