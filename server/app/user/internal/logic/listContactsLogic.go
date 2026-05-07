package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListContactsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListContactsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactsLogic {
	return &ListContactsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListContactsLogic) ListContacts(in *user.ListContactsRequest) (*user.ListContactsResponse, error) {
	userId, _ := strconv.ParseInt(in.UserId, 10, 64)
	rows, err := l.svcCtx.ContactDao.List(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	var contacts []*user.ContactInfo
	for _, r := range rows {
		contacts = append(contacts, &user.ContactInfo{
			UserId: strconv.FormatInt(r.ContactId, 10),
			Nickname: r.Nickname, Avatar: r.Avatar,
			Remark: r.Remark, Tag: r.Tag, IsBlocked: r.IsBlocked,
		})
	}
	return &user.ListContactsResponse{Contacts: contacts}, nil
}
