// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"auth4x/internal/logic"
	"auth4x/internal/svc"
	"auth4x/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SignUpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignUpReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSignUpLogic(r.Context(), svcCtx)
		resp, err := l.SignUp(&req)
		if err != nil {
			writeError(r, w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
