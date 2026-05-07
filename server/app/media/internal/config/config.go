package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Minio struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
	}
}
