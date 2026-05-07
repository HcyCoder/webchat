package handler

import (
	"net/http"
	"github.com/team/webchat-server/app/gateway/internal/logic"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetMessagesLogic(r.Context(), svcCtx)
		resp, err := l.GetMessages(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
