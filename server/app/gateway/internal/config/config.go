package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpc  zrpc.RpcClientConf
	ChatRpc  zrpc.RpcClientConf
	GroupRpc zrpc.RpcClientConf
	MediaRpc zrpc.RpcClientConf
	Cache    struct {
		Host string
	}
}
