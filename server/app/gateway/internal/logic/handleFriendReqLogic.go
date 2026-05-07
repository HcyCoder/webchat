package logic

import (
	"context"
	"net/http"

	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFriendReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleFriendReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFriendReqLogic {
	return &HandleFriendReqLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleFriendReqLogic) HandleFriendReq(r *http.Request, req *types.HandleFriendReqRequest) (*types.Empty, error) {
	requestID := extractPathParam(r.URL.Path, "request/")
	if req.Action == "accept" {
		_, err := l.svcCtx.UserRpc.AcceptFriendRequest(l.ctx, &userClient.AcceptFriendRequestRequest{RequestId: requestID})
		return &types.Empty{}, err
	}
	_, err := l.svcCtx.UserRpc.RejectFriendRequest(l.ctx, &userClient.RejectFriendRequestRequest{RequestId: requestID})
	return &types.Empty{}, err
}

