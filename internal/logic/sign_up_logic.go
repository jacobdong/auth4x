// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"net/http"

	"auth4x/internal/errs"
	"auth4x/internal/store"
	"auth4x/internal/svc"
	"auth4x/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignUpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignUpLogic {
	return &SignUpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignUpLogic) SignUp(req *types.SignUpReq) (resp *types.SignUpResp, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errs.New(http.StatusBadRequest, "username and password are required")
	}

	user, err := l.svcCtx.Store.CreateUser(req.Username, req.Password)
	if err != nil {
		if err == store.ErrUserExists {
			return nil, errs.New(http.StatusConflict, "user already exists")
		}
		return nil, err
	}

	return &types.SignUpResp{UserId: user.ID}, nil
}
