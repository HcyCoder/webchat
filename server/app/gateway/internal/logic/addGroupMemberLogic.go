package logic

import (
	"context"
	"net/http"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/group/groupClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddGroupMemberLogic) AddGroupMember(r *http.Request, req *types.AddMemberRequest) (*types.Empty, error) {
	gid := extractPathParam(r.URL.Path, "groups/")
	_, err := l.svcCtx.GroupRpc.AddMember(l.ctx, &groupClient.AddMemberRequest{GroupId: gid, UserIds: req.UserIds})
	return &types.Empty{}, err
}

