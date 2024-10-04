package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	xhttp "github.com/zeromicro/x/http"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shortURL/internal/logic"
	"shortURL/internal/svc"
	"shortURL/internal/types"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowReq
		if err := httpx.Parse(r, &req); err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		//validator参数校验
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(&req)
		if err != nil {
			logx.Errorw("validator 校验参数失败", logx.LogField{Key: "err", Value: err})
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			//httpx.ErrorCtx(r.Context(), w, err)
		} else {
			//xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
			//httpx.OkJsonCtx(r.Context(), w, resp)
			http.Redirect(w, r, resp.LongURL, http.StatusFound) // 重定向
		}
	}
}
