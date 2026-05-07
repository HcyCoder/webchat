package svc

import (
	"github.com/team/webchat-server/app/user/internal/config"
	"github.com/team/webchat-server/app/user/internal/dao"
	"github.com/team/webchat-server/common/mysql"
	"github.com/team/webchat-server/common/token"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config       config.Config
	UserDao      *dao.UserDao
	ContactDao   *dao.ContactDao
	TokenManager *token.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := mysql.New(c.Mysql.DataSource)
	rdb := redis.NewClient(&redis.Options{Addr: c.Redis.Addr})
	return &ServiceContext{
		Config:       c,
		UserDao:      dao.NewUserDao(conn),
		ContactDao:   dao.NewContactDao(conn),
		TokenManager: token.NewManager(rdb),
	}
}
