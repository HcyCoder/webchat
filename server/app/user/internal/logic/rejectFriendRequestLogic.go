package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type RejectFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectFriendRequestLogic {
	return &RejectFriendRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RejectFriendRequestLogic) RejectFriendRequest(in *user.RejectFriendRequestRequest) (*user.Empty, error) {
	requestId, _ := strconv.ParseInt(in.RequestId, 10, 64)
	err := l.svcCtx.ContactDao.RejectRequest(l.ctx, requestId)
	return &user.Empty{}, err
}
