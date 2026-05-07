package logic

import (
	"context"
	"github.com/team/webchat-server/app/chat/chatClient"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMsgRequest) (*types.Message, error) {
	uid := middleware.GetUserID(l.ctx)
	m, err := l.svcCtx.ChatRpc.SendMessage(l.ctx, &chatClient.SendMessageRequest{
		ChatType: req.ChatType, FromUser: uid, ToId: req.ToId, MsgType: req.MsgType, Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &types.Message{
		Id: m.Id, ChatType: m.ChatType, FromUser: m.FromUser, ToId: m.ToId,
		MsgType: m.MsgType, Content: m.Content, IsRecalled: m.IsRecalled, CreatedAt: m.CreatedAt,
	}, nil
}
