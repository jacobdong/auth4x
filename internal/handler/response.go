package handler

import (
	"net/http"

	"auth4x/internal/errs"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func writeError(r *http.Request, w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*errs.APIError); ok {
		httpx.WriteJson(w, apiErr.Status, map[string]string{"error": apiErr.Message})
		return
	}

	httpx.ErrorCtx(r.Context(), w, err)
}
