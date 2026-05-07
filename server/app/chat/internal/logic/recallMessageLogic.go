package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RecallMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RecallMessageLogic) RecallMessage(in *chat.RecallMessageRequest) (*chat.Empty, error) {
	msgId, _ := strconv.ParseInt(in.MsgId, 10, 64)
	userId, _ := strconv.ParseInt(in.UserId, 10, 64)
	err := l.svcCtx.MessageDao.Recall(l.ctx, msgId, userId)
	return &chat.Empty{}, err
}
