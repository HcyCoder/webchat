package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddContactLogic {
	return &AddContactLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddContactLogic) AddContact(in *user.AddContactRequest) (*user.Empty, error) {
	fromUser, _ := strconv.ParseInt(in.FromUser, 10, 64)
	toUser, _ := strconv.ParseInt(in.ToUser, 10, 64)
	err := l.svcCtx.ContactDao.AddRequest(l.ctx, fromUser, toUser, in.Message)
	return &user.Empty{}, err
}
