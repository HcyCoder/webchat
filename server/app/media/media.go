package main

import (
	"flag"

	"github.com/team/webchat-server/app/media/internal/config"
	"github.com/team/webchat-server/app/media/internal/server"
	"github.com/team/webchat-server/app/media/internal/svc"
	"github.com/team/webchat-server/app/media/internal/media"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/media.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		media.RegisterMediaServer(grpcServer, server.NewMediaServer(ctx))
	})
	defer s.Stop()
	s.Start()
}
