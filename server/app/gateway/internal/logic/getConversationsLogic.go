package logic

import (
	"context"
	"github.com/team/webchat-server/app/chat/chatClient"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetConversationsLogic) GetConversations() (*types.GetConversationsResponse, error) {
	uid := middleware.GetUserID(l.ctx)
	resp, err := l.svcCtx.ChatRpc.GetConversations(l.ctx, &chatClient.GetConversationsRequest{UserId: uid})
	if err != nil {
		return nil, err
	}
	var convs []types.Conversation
	for _, c := range resp.Conversations {
		convs = append(convs, types.Conversation{
			Id: c.Id, ChatType: c.ChatType, TargetId: c.TargetId,
			TargetName: c.TargetName, TargetAvatar: c.TargetAvatar,
			UnreadCount: int(c.UnreadCount), IsPinned: c.IsPinned, IsMuted: c.IsMuted, UpdatedAt: c.UpdatedAt,
		})
	}
	return &types.GetConversationsResponse{Conversations: convs}, nil
}
