package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberLogic {
	return &AddMemberLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddMemberLogic) AddMember(in *group.AddMemberRequest) (*group.Empty, error) {
	gid, _ := strconv.ParseInt(in.GroupId, 10, 64)
	var ids []int64
	for _, s := range in.UserIds {
		id, _ := strconv.ParseInt(s, 10, 64)
		ids = append(ids, id)
	}
	err := l.svcCtx.GroupDao.AddMembers(l.ctx, gid, ids)
	return &group.Empty{}, err
}
