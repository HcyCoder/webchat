package main

import (
	"flag"

	"github.com/team/webchat-server/app/group/internal/config"
	"github.com/team/webchat-server/app/group/internal/server"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/group.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		group.RegisterGroupServer(grpcServer, server.NewGroupServer(ctx))
	})
	defer s.Stop()
	s.Start()
}
