package handler

import (
	"net/http"
	"github.com/team/webchat-server/app/gateway/internal/logic"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetGroupLogic(r.Context(), svcCtx)
		resp, err := l.GetGroup(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
