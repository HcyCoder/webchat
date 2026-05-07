package logic

import (
	"context"
	"net/http"
	"github.com/team/webchat-server/app/chat/chatClient"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMessagesLogic) GetMessages(r *http.Request) (*types.GetMessagesResponse, error) {
	convID := extractPathParam(r.URL.Path, "messages/")
	resp, err := l.svcCtx.ChatRpc.GetMessages(l.ctx, &chatClient.GetMessagesRequest{ConvId: convID, Page: 1, PageSize: 50})
	if err != nil {
		return nil, err
	}
	var msgs []types.Message
	for _, m := range resp.Messages {
		msgs = append(msgs, types.Message{
			Id: m.Id, ChatType: m.ChatType, FromUser: m.FromUser, ToId: m.ToId,
			MsgType: m.MsgType, Content: m.Content, IsRecalled: m.IsRecalled, CreatedAt: m.CreatedAt,
		})
	}
	return &types.GetMessagesResponse{Messages: msgs}, nil
}

