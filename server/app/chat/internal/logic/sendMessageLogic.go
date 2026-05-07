package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/team/webchat-server/app/chat/model"
	"github.com/team/webchat-server/app/chat/ws"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *chat.SendMessageRequest) (*chat.Message, error) {
	fromUser, _ := strconv.ParseInt(in.FromUser, 10, 64)
	toId, _ := strconv.ParseInt(in.ToId, 10, 64)
	msg := &model.Message{ChatType: in.ChatType, FromUser: fromUser, ToId: toId, MsgType: in.MsgType, Content: in.Content}
	id, err := l.svcCtx.MessageDao.Insert(l.ctx, msg)
	if err != nil {
		return nil, err
	}
	l.svcCtx.ConversationDao.Upsert(l.ctx, fromUser, toId, in.ChatType, id)
	l.svcCtx.ConversationDao.Upsert(l.ctx, toId, fromUser, in.ChatType, id)
	payload, _ := json.Marshal(msg)
	l.svcCtx.Hub.SendToUser(toId, &ws.WsMessage{Type: "chat.message", Seq: id, Data: payload})
	return &chat.Message{
		Id: strconv.FormatInt(id, 10), ChatType: msg.ChatType,
		FromUser: in.FromUser, ToId: in.ToId,
		MsgType: msg.MsgType, Content: msg.Content, CreatedAt: msg.CreatedAt,
	}, nil
}
