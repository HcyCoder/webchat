package svc

import (
	"github.com/team/webchat-server/app/gateway/internal/config"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/team/webchat-server/app/chat/chatClient"
	"github.com/team/webchat-server/app/group/groupClient"
	"github.com/team/webchat-server/app/media/mediaClient"
	"github.com/team/webchat-server/common/token"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	UserRpc      userClient.User
	ChatRpc      chatClient.Chat
	GroupRpc     groupClient.Group
	MediaRpc     mediaClient.Media
	TokenManager *token.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{Addr: c.Cache.Host})
	return &ServiceContext{
		Config:       c,
		UserRpc:      userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ChatRpc:      chatClient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		GroupRpc:     groupClient.NewGroup(zrpc.MustNewClient(c.GroupRpc)),
		MediaRpc:     mediaClient.NewMedia(zrpc.MustNewClient(c.MediaRpc)),
		TokenManager: token.NewManager(rdb),
	}
}
