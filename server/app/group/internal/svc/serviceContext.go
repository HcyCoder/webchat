package svc

import (
	"github.com/team/webchat-server/app/group/internal/config"
	"github.com/team/webchat-server/app/group/internal/dao"
	"github.com/team/webchat-server/common/mysql"
)

type ServiceContext struct {
	Config   config.Config
	GroupDao *dao.GroupDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := mysql.New(c.Mysql.DataSource)
	return &ServiceContext{Config: c, GroupDao: dao.NewGroupDao(conn)}
}
