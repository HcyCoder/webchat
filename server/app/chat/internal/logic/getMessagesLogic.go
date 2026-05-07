package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetMessagesLogic) GetMessages(in *chat.GetMessagesRequest) (*chat.GetMessagesResponse, error) {
	convId, _ := strconv.ParseInt(in.ConvId, 10, 64)
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	msgs, err := l.svcCtx.MessageDao.FindByConv(l.ctx, convId, "single", page, int(in.PageSize))
	if err != nil {
		return nil, err
	}
	var pbMsgs []*chat.Message
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, &chat.Message{
			Id: strconv.FormatInt(m.Id, 10), ChatType: m.ChatType,
			FromUser: strconv.FormatInt(m.FromUser, 10), ToId: strconv.FormatInt(m.ToId, 10),
			MsgType: m.MsgType, Content: m.Content, IsRecalled: m.IsRecalled, CreatedAt: m.CreatedAt,
		})
	}
	return &chat.GetMessagesResponse{Messages: pbMsgs}, nil
}
