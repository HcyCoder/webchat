package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptFriendRequestLogic {
	return &AcceptFriendRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AcceptFriendRequestLogic) AcceptFriendRequest(in *user.AcceptFriendRequestRequest) (*user.Empty, error) {
	requestId, _ := strconv.ParseInt(in.RequestId, 10, 64)
	_, _, err := l.svcCtx.ContactDao.AcceptRequest(l.ctx, requestId)
	return &user.Empty{}, err
}
