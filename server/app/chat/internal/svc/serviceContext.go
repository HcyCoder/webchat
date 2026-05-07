package svc

import (
	"github.com/team/webchat-server/app/chat/internal/config"
	"github.com/team/webchat-server/app/chat/internal/dao"
	"github.com/team/webchat-server/app/chat/ws"
	"github.com/team/webchat-server/common/mysql"
	"github.com/team/webchat-server/common/token"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config         config.Config
	MessageDao     *dao.MessageDao
	ConversationDao *dao.ConversationDao
	Hub            *ws.Hub
	TokenManager   *token.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := mysql.New(c.Mysql.DataSource)
	rdb := redis.NewClient(&redis.Options{Addr: c.Redis.Addr})
	return &ServiceContext{
		Config:          c,
		MessageDao:      dao.NewMessageDao(conn),
		ConversationDao: dao.NewConversationDao(conn),
		Hub:             ws.NewHub(),
		TokenManager:    token.NewManager(rdb),
	}
}
