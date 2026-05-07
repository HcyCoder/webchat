package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddContactLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddContactLogic {
	return &AddContactLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddContactLogic) AddContact(req *types.AddContactRequest) (*types.Empty, error) {
	uid := middleware.GetUserID(l.ctx)
	_, err := l.svcCtx.UserRpc.AddContact(l.ctx, &userClient.AddContactRequest{
		FromUser: uid, ToUser: req.ToUser, Message: req.Message,
	})
	return &types.Empty{}, err
}
