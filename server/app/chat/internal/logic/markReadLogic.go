package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkReadLogic) MarkRead(in *chat.MarkReadRequest) (*chat.Empty, error) {
	userId, _ := strconv.ParseInt(in.UserId, 10, 64)
	convId, _ := strconv.ParseInt(in.ConvId, 10, 64)
	err := l.svcCtx.ConversationDao.ClearUnread(l.ctx, convId, userId)
	return &chat.Empty{}, err
}
