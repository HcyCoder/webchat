package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMembersLogic {
	return &ListMembersLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListMembersLogic) ListMembers(in *group.ListMembersRequest) (*group.ListMembersResponse, error) {
	gid, _ := strconv.ParseInt(in.GroupId, 10, 64)
	rows, err := l.svcCtx.GroupDao.ListMembers(l.ctx, gid)
	if err != nil {
		return nil, err
	}
	var members []*group.GroupMember
	for _, r := range rows {
		members = append(members, &group.GroupMember{
			UserId: strconv.FormatInt(r.UserId, 10), Role: r.Role,
			Alias: r.Alias, IsMuted: r.IsMuted, JoinedAt: r.JoinedAt,
		})
	}
	return &group.ListMembersResponse{Members: members}, nil
}
