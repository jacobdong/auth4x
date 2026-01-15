// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package users

import (
	"net/http"

	"auth4x/internal/logic/users"
	"auth4x/internal/svc"
	"auth4x/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserMeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewUserMeLogic(r.Context(), svcCtx)
		resp, err := l.UserMe(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
