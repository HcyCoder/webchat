package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListContactsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListContactsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactsLogic {
	return &ListContactsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListContactsLogic) ListContacts() (*types.ListContactsResponse, error) {
	uid := middleware.GetUserID(l.ctx)
	resp, err := l.svcCtx.UserRpc.ListContacts(l.ctx, &userClient.ListContactsRequest{UserId: uid})
	if err != nil {
		return nil, err
	}
	var contacts []types.ContactInfo
	for _, c := range resp.Contacts {
		contacts = append(contacts, types.ContactInfo{
			UserId: c.UserId, Nickname: c.Nickname, Avatar: c.Avatar,
			Remark: c.Remark, Tag: c.Tag, IsBlocked: c.IsBlocked,
		})
	}
	return &types.ListContactsResponse{Contacts: contacts}, nil
}
