package main

import (
	"flag"

	"github.com/team/webchat-server/app/gateway/internal/config"
	"github.com/team/webchat-server/app/gateway/internal/handler"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	handler.RegisterHandlers(server, ctx)
	server.Start()
}
