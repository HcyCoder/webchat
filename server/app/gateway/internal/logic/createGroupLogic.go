package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/group/groupClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateGroupLogic) CreateGroup(req *types.CreateGroupRequest) (*types.GroupInfo, error) {
	uid := middleware.GetUserID(l.ctx)
	g, err := l.svcCtx.GroupRpc.CreateGroup(l.ctx, &groupClient.CreateGroupRequest{
		Name: req.Name, OwnerId: uid, MemberIds: req.MemberIds,
	})
	if err != nil {
		return nil, err
	}
	return &types.GroupInfo{
		Id: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId,
		Announcement: g.Announcement, MemberCount: int(g.MemberCount),
	}, nil
}
