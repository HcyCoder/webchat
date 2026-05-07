package logic

import (
	"context"
	"net/http"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/group/groupClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupLogic {
	return &GetGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupLogic) GetGroup(r *http.Request) (*types.GroupInfo, error) {
	gid := extractPathParam(r.URL.Path, "groups/")
	g, err := l.svcCtx.GroupRpc.GetGroup(l.ctx, &groupClient.GetGroupRequest{GroupId: gid})
	if err != nil {
		return nil, err
	}
	return &types.GroupInfo{
		Id: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId,
		Announcement: g.Announcement, MemberCount: int(g.MemberCount),
	}, nil
}

