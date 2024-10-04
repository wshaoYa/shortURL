package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"shortURL/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shortURL/internal/logic"
	"shortURL/internal/svc"

	xhttp "github.com/zeromicro/x/http" // 自定义 resp响应格式
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertReq
		if err := httpx.Parse(r, &req); err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		//validator参数校验
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(&req)
		if err != nil {
			logx.Infow("validator 校验参数失败", logx.Field("err", err))
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		//业务逻辑
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			//httpx.OkJsonCtx(r.Context(), w, resp)
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
