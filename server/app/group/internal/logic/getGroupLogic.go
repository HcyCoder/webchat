package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupLogic {
	return &GetGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetGroupLogic) GetGroup(in *group.GetGroupRequest) (*group.GroupInfo, error) {
	gid, _ := strconv.ParseInt(in.GroupId, 10, 64)
	g, err := l.svcCtx.GroupDao.FindById(l.ctx, gid)
	if err != nil {
		return nil, err
	}
	return &group.GroupInfo{
		Id: strconv.FormatInt(g.Id, 10), Name: g.Name, Avatar: g.Avatar,
		OwnerId: strconv.FormatInt(g.OwnerId, 10), Announcement: g.Announcement,
		MemberCount: g.MemberCount, MaxMembers: g.MaxMembers, CreatedAt: g.CreatedAt,
	}, nil
}
