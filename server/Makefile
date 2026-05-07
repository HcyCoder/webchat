.PHONY: api rpc run-gateway run-user run-chat run-group run-media

api:
	goctl api go -api app/gateway/gateway.api -dir app/gateway -style goZero

rpc:
	cd app/user    && goctl rpc protoc user.proto    --go_out=./internal --go-grpc_out=./internal --zrpc_out=. --style goZero
	cd app/chat    && goctl rpc protoc chat.proto    --go_out=./internal --go-grpc_out=./internal --zrpc_out=. --style goZero
	cd app/group   && goctl rpc protoc group.proto   --go_out=./internal --go-grpc_out=./internal --zrpc_out=. --style goZero
	cd app/media   && goctl rpc protoc media.proto   --go_out=./internal --go-grpc_out=./internal --zrpc_out=. --style goZero

run-gateway:
	go run ./app/gateway/gateway.go -f app/gateway/etc/gateway.yaml

run-user:
	go run ./app/user/user.go -f app/user/etc/user.yaml

run-chat:
	go run ./app/chat/chat.go -f app/chat/etc/chat.yaml

run-group:
	go run ./app/group/group.go -f app/group/etc/group.yaml

run-media:
	go run ./app/media/media.go -f app/media/etc/media.yaml
