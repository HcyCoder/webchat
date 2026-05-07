package main

import (
	"flag"

	"github.com/team/webchat-server/app/chat/internal/config"
	"github.com/team/webchat-server/app/chat/internal/server"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/team/webchat-server/app/chat/internal/chat"
	"github.com/team/webchat-server/app/chat/ws"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/chat.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	go func() {
		s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
			chat.RegisterChatServer(grpcServer, server.NewChatServer(ctx))
		})
		defer s.Stop()
		s.Start()
	}()

	if err := ws.NewServer(ctx.Hub, ctx.TokenManager).Listen(c.WebSocket.ListenOn); err != nil {
		panic(err)
	}
}
