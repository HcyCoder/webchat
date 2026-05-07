package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMemberLogic {
	return &RemoveMemberLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RemoveMemberLogic) RemoveMember(in *group.RemoveMemberRequest) (*group.Empty, error) {
	gid, _ := strconv.ParseInt(in.GroupId, 10, 64)
	uid, _ := strconv.ParseInt(in.UserId, 10, 64)
	err := l.svcCtx.GroupDao.RemoveMember(l.ctx, gid, uid)
	return &group.Empty{}, err
}
