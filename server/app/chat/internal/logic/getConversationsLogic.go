package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetConversationsLogic) GetConversations(in *chat.GetConversationsRequest) (*chat.GetConversationsResponse, error) {
	userId, _ := strconv.ParseInt(in.UserId, 10, 64)
	rows, err := l.svcCtx.ConversationDao.ListByUser(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	var convs []*chat.Conversation
	for _, r := range rows {
		convs = append(convs, &chat.Conversation{
			Id: strconv.FormatInt(r.Id, 10), ChatType: r.ChatType,
			TargetId: strconv.FormatInt(r.TargetId, 10),
			UnreadCount: uint32(r.UnreadCnt), IsPinned: r.IsPinned, IsMuted: r.IsMuted, UpdatedAt: r.UpdatedAt,
		})
	}
	return &chat.GetConversationsResponse{Conversations: convs}, nil
}
