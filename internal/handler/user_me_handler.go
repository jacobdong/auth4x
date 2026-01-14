// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"auth4x/internal/logic"
	"auth4x/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserMeLogic(r.Context(), svcCtx, r)
		resp, err := l.UserMe()
		if err != nil {
			writeError(r, w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
