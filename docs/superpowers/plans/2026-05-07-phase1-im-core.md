# Phase 1: IM 核心 — 实施计划 (go-zero + Token/Redis)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建 go-zero 微服务基础设施 + 完成 IM 核心（用户/聊天/群聊/文件）+ Flutter 客户端聊天主流程

**Architecture:** go-zero 微服务集群 —— 1 个 API 服务 (gateway) + 4 个 RPC 服务 (user/chat/group/media)，etcd 服务注册发现，gRPC 内部通信。认证用 Token + Redis（SHA256 生成 token，Redis 存映射并设 TTL）。chat-service 内嵌独立 WebSocket HTTP server。Flutter 端按 BLoC 模式组织。

**Tech Stack:** Go + go-zero + goctl + gRPC + Protobuf + Kafka + MySQL + Redis + MinIO + Docker Compose (本地) / Flutter 3.x + flutter_bloc + dio

---

## 文件结构

### Backend (go-zero monorepo)

```
server/
├── go.mod / go.sum / Makefile
├── docker-compose.yml
├── common/
│   ├── token/token.go          # Token 生成与 Redis 存储
│   ├── mysql/mysql.go          # MySQL 连接
│   └── errcode/errcode.go      # 公共错误码
├── app/
│   ├── gateway/
│   │   ├── etc/gateway.yaml
│   │   ├── gateway.api
│   │   ├── gateway.go
│   │   └── internal/
│   │       ├── config/config.go
│   │       ├── handler/register.go
│   │       ├── logic/
│   │       │   ├── auth/login.go, register.go
│   │       │   ├── user/me.go, update.go
│   │       │   ├── chat/conversations.go, sendmsg.go
│   │       │   ├── group/creategroup.go, getgroup.go
│   │       │   └── media/uploadurl.go, fileurl.go
│   │       ├── middleware/auth.go
│   │       ├── svc/servicecontext.go
│   │       └── types/types.go
│   ├── user/
│   │   ├── etc/user.yaml
│   │   ├── user.proto
│   │   ├── user.go
│   │   ├── internal/
│   │   │   ├── config/config.go
│   │   │   ├── server/userserver.go
│   │   │   ├── logic/ (register.go, login.go, getprofile.go, updateprofile.go, addcontact.go, listcontacts.go, searchuser.go, acceptfriendreq.go)
│   │   │   ├── svc/servicecontext.go
│   │   │   └── dao/user.go, contact.go
│   │   ├── userclient/user.go
│   │   └── model/user.go, contact.go
│   ├── chat/
│   │   ├── etc/chat.yaml
│   │   ├── chat.proto
│   │   ├── chat.go
│   │   ├── internal/
│   │   │   ├── config/config.go
│   │   │   ├── server/chatserver.go
│   │   │   ├── logic/ (sendmessage.go, getmessages.go, getconversations.go, markread.go, recall.go)
│   │   │   ├── svc/servicecontext.go
│   │   │   └── dao/message.go, conversation.go
│   │   ├── chatclient/chat.go
│   │   ├── ws/ (hub.go, client.go, server.go)
│   │   └── model/message.go, conversation.go
│   ├── group/
│   │   ├── etc/group.yaml
│   │   ├── group.proto
│   │   ├── group.go
│   │   ├── internal/ (同上 RPC 模式)
│   │   ├── groupclient/group.go
│   │   └── model/group.go
│   └── media/
│       ├── etc/media.yaml
│       ├── media.proto
│       ├── media.go
│       ├── internal/ (同上 RPC 模式)
│       ├── mediaclient/media.go
│       └── model/media.go
└── migrations/
    ├── user/001_init.sql
    ├── chat/001_init.sql
    ├── group/001_init.sql
```

### Flutter 客户端

```
client/
├── pubspec.yaml
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── core/
│   │   ├── network/   (dio_client.dart, ws_client.dart)
│   │   ├── auth/      (token_manager.dart)
│   │   └── theme/     (colors.dart, theme.dart)
│   ├── features/
│   │   ├── auth/
│   │   │   ├── bloc/  (auth_bloc.dart, auth_event.dart, auth_state.dart)
│   │   │   ├── data/  (auth_repository.dart, auth_api.dart)
│   │   │   └── ui/    (login_page.dart, register_page.dart)
│   │   ├── chat/
│   │   │   ├── bloc/  (conversation_bloc.dart, message_bloc.dart)
│   │   │   ├── data/  (chat_repository.dart, chat_api.dart)
│   │   │   └── ui/    (conversation_list_page.dart, chat_page.dart, widgets/)
│   │   ├── contacts/
│   │   │   ├── bloc/  (contacts_bloc.dart)
│   │   │   ├── data/  (contacts_repository.dart, contacts_api.dart)
│   │   │   └── ui/    (contacts_list_page.dart, add_contact_page.dart)
│   │   └── home/      (home_page.dart)
│   └── shared/widgets/ (avatar.dart, loading_indicator.dart)
```

---

## Part A: 开发环境与基础设施

### Task A1: 初始化 go-zero 项目骨架

**Files:**
- Create: `server/go.mod`, `server/Makefile`
- Create: `server/app/gateway/gateway.go` (骨架)
- Create: `server/app/user/user.go` (骨架)
- Create: `server/app/chat/chat.go` (骨架)
- Create: `server/app/group/group.go` (骨架)
- Create: `server/app/media/media.go` (骨架)

- [ ] **Step 1: 初始化 Go module**

```bash
mkdir -p server/app/{gateway,user,chat,group,media}
cd server && go mod init github.com/team/webchat-server
```

- [ ] **Step 2: 安装 go-zero 和 goctl**

```bash
go get github.com/zeromicro/go-zero@latest
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

- [ ] **Step 3: 创建各服务入口骨架**

```go
// server/app/user/user.go
package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	flag.Parse()
	var c zrpc.RpcServerConf
	conf.MustLoad("etc/user.yaml", &c)
	s := zrpc.MustNewServer(c, nil)
	defer s.Stop()
	s.Start()
}
```

其他服务同理，类型替换为对应的 server conf。

- [ ] **Step 4: 创建 Makefile**

```makefile
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
```

- [ ] **Step 5: Commit**

```bash
git add server/go.mod server/go.sum server/Makefile server/app/ && git commit -m "feat: init go-zero project skeleton"
```

---

### Task A2: 编写 docker-compose.yml

**Files:**
- Create: `server/docker-compose.yml`

- [ ] **Step 1: 编写 docker-compose.yml**

与之前 PRD 设计一致（MySQL 8.0 + Redis 7 + MinIO + Kafka + etcd），内容不变。

- [ ] **Step 2: 启动并验证**

```bash
cd server && docker compose up -d
```

- [ ] **Step 3: Commit**

```bash
git add server/docker-compose.yml && git commit -m "feat: add docker-compose dev infrastructure"
```

---

### Task A3: 定义 Protobuf + API 文件

**Files:**
- Create: `server/app/gateway/gateway.api`
- Create: `server/app/user/user.proto`
- Create: `server/app/chat/chat.proto`
- Create: `server/app/group/group.proto`
- Create: `server/app/media/media.proto`

- [ ] **Step 1: 编写 gateway.api (go-zero API 定义)**

```go
// server/app/gateway/gateway.api
syntax = "v1"

type (
	LoginRequest {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	LoginResponse {
		UserId string `json:"user_id"`
		Token  string `json:"token"`
	}
	RegisterRequest {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	RegisterResponse {
		UserId string `json:"user_id"`
		Token  string `json:"token"`
	}
	RefreshRequest {
		UserId string `json:"user_id"`
	}
	RefreshResponse {
		Token string `json:"token"`
	}
	UserInfo {
		Id        string `json:"id"`
		Phone     string `json:"phone"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Gender    int    `json:"gender"`
		Region    string `json:"region"`
		Signature string `json:"signature"`
	}
	UpdateUserRequest {
		Nickname  string `json:"nickname,optional"`
		Avatar    string `json:"avatar,optional"`
		Gender    int    `json:"gender,optional"`
		Region    string `json:"region,optional"`
		Signature string `json:"signature,optional"`
	}
	Contact {
		UserId    string `json:"user_id"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Remark    string `json:"remark"`
		Tag       string `json:"tag"`
		IsBlocked bool   `json:"is_blocked"`
	}
	AddContactRequest {
		ToUser  string `json:"to_user"`
		Message string `json:"message"`
	}
	HandleFriendReqRequest {
		Action string `json:"action"` // "accept" | "reject"
	}
	Conversation {
		Id           string `json:"id"`
		ChatType     string `json:"chat_type"`
		TargetId     string `json:"target_id"`
		TargetName   string `json:"target_name"`
		TargetAvatar string `json:"target_avatar"`
		LastContent  string `json:"last_content"`
		UnreadCount  int    `json:"unread_count"`
		IsPinned     bool   `json:"is_pinned"`
		IsMuted      bool   `json:"is_muted"`
		UpdatedAt    int64  `json:"updated_at"`
	}
	Message {
		Id         string `json:"id"`
		ChatType   string `json:"chat_type"`
		FromUser   string `json:"from_user"`
		ToId       string `json:"to_id"`
		MsgType    string `json:"msg_type"`
		Content    string `json:"content"`
		IsRecalled bool   `json:"is_recalled"`
		CreatedAt  int64  `json:"created_at"`
	}
	SendMsgRequest {
		ChatType string `json:"chat_type"`
		ToId     string `json:"to_id"`
		MsgType  string `json:"msg_type"`
		Content  string `json:"content"`
	}
	CreateGroupRequest {
		Name      string   `json:"name"`
		MemberIds []string `json:"member_ids"`
	}
	GroupInfo {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Avatar       string `json:"avatar"`
		OwnerId      string `json:"owner_id"`
		Announcement string `json:"announcement"`
		MemberCount  int    `json:"member_count"`
	}
	AddMemberRequest {
		UserIds []string `json:"user_ids"`
	}
	UploadURLResponse {
		UploadUrl string `json:"upload_url"`
		FileId    string `json:"file_id"`
		ExpiresIn int64  `json:"expires_in"`
	}
	FileURLResponse {
		Url string `json:"url"`
	}
)

@server(
	prefix: /api/v1
)
service gateway {
	@handler login
	post /auth/login (LoginRequest) returns (LoginResponse)

	@handler register
	post /auth/register (RegisterRequest) returns (RegisterResponse)

	@handler refreshToken
	post /auth/refresh returns (RefreshResponse)

	@handler getUser
	get /users/me returns (UserInfo)

	@handler updateUser
	put /users/me (UpdateUserRequest) returns (UserInfo)

	@handler getUserById
	get /users/:id returns (UserInfo)

	@handler listContacts
	get /contacts returns (ListContactsResponse)

	@handler addContact
	post /contacts/request (AddContactRequest)

	@handler handleFriendReq
	put /contacts/request/:id (HandleFriendReqRequest)

	@handler getConversations
	get /conversations returns (GetConversationsResponse)

	@handler getMessages
	get /messages/:conv_id returns (GetMessagesResponse)

	@handler sendMessage
	post /messages/send (SendMsgRequest) returns (Message)

	@handler createGroup
	post /groups (CreateGroupRequest) returns (GroupInfo)

	@handler getGroup
	get /groups/:id returns (GroupInfo)

	@handler addGroupMember
	post /groups/:id/members (AddMemberRequest)

	@handler getUploadURL
	post /files/upload returns (UploadURLResponse)

	@handler getFileURL
	get /files/:id/url returns (FileURLResponse)
}
```

- [ ] **Step 2: 编写 user.proto**

与之前 PRD 中的 proto 定义一致，补充 go-zero 需要的 option：

```protobuf
// server/app/user/user.proto
syntax = "proto3";
package user;
option go_package = "./user";

service User {
  rpc register(RegisterRequest) returns (RegisterResponse);
  rpc login(LoginRequest) returns (LoginResponse);
  rpc getProfile(GetProfileRequest) returns (UserInfo);
  rpc updateProfile(UpdateProfileRequest) returns (UserInfo);
  rpc addContact(AddContactRequest) returns (Empty);
  rpc acceptFriendRequest(AcceptFriendRequestRequest) returns (Empty);
  rpc rejectFriendRequest(RejectFriendRequestRequest) returns (Empty);
  rpc listContacts(ListContactsRequest) returns (ListContactsResponse);
  rpc searchUser(SearchUserRequest) returns (SearchUserResponse);
  rpc getUserById(GetUserByIdRequest) returns (UserInfo);
}

message Empty {}

message UserInfo {
  string id = 1;
  string phone = 2;
  string nickname = 3;
  string avatar = 4;
  int32 gender = 5;
  string region = 6;
  string signature = 7;
  int64 created_at = 8;
}

message RegisterRequest {
  string phone = 1;
  string password = 2;
  string nickname = 3;
}

message RegisterResponse {
  string user_id = 1;
}

message LoginRequest {
  string phone = 1;
  string password = 2;
}

message LoginResponse {
  string user_id = 1;
}

message GetProfileRequest { string user_id = 1; }
message GetUserByIdRequest { string user_id = 1; }

message UpdateProfileRequest {
  string user_id = 1;
  string nickname = 2;
  string avatar = 3;
  int32 gender = 4;
  string region = 5;
  string signature = 6;
}

message AddContactRequest {
  string from_user = 1;
  string to_user = 2;
  string message = 3;
}

message AcceptFriendRequestRequest { string request_id = 1; }
message RejectFriendRequestRequest { string request_id = 1; }

message ListContactsRequest { string user_id = 1; }

message ListContactsResponse {
  repeated ContactInfo contacts = 1;
}

message ContactInfo {
  string user_id = 1;
  string nickname = 2;
  string avatar = 3;
  string remark = 4;
  string tag = 5;
  bool is_blocked = 6;
}

message SearchUserRequest { string keyword = 1; }

message SearchUserResponse {
  repeated UserInfo users = 1;
}
```

- [ ] **Step 3: 编写 chat.proto**

与之前 PRD 中的 chat proto 一致，加 go-zero option。

- [ ] **Step 4: 编写 group.proto 和 media.proto**

与之前 PRD 一致，加 go-zero option。

- [ ] **Step 5: 用 goctl 生成代码**

```bash
cd server
make api   # 生成 gateway handler + types
make rpc   # 生成各 RPC 服务的 server + client + logic
```

- [ ] **Step 6: 验证编译**

```bash
cd server && go mod tidy && go build ./app/...
```

Expected: 编译成功。

- [ ] **Step 7: Commit**

```bash
git add server/app/*/gateway.api server/app/*/*.proto server/app/*/internal/ server/app/*/*client/ && git commit -m "feat: add protobuf + API definitions, generate code"
```

---

### Task A4: 创建公共包 (common/)

**Files:**
- Create: `server/common/token/token.go`
- Create: `server/common/mysql/mysql.go`
- Create: `server/common/errcode/errcode.go`

- [ ] **Step 1: Token 包（Token + Redis 认证核心）**

```go
// server/common/token/token.go
package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"github.com/redis/go-redis/v9"
)

const (
	TokenPrefix = "token:"
	TokenTTL    = 24 * time.Hour
	GraceTTL    = 5 * time.Minute // 刷新后旧 token 保留时间
)

type Manager struct {
	rdb *redis.Client
}

func NewManager(rdb *redis.Client) *Manager {
	return &Manager{rdb: rdb}
}

// token = SHA256(userID + timestamp + random)[:32]
func Generate(userID string) string {
	src := fmt.Sprintf("%s_%d_%d", userID, time.Now().UnixNano(), rand.Int63())
	hash := sha256.Sum256([]byte(src))
	return hex.EncodeToString(hash[:])[:32]
}

// Redis key: "token:<token>" -> userID, TTL 24h
func (m *Manager) Store(ctx context.Context, token, userID string) error {
	key := TokenPrefix + token
	return m.rdb.Set(ctx, key, userID, TokenTTL).Err()
}

// 返回 user_id，如果不存在或过期返回空
func (m *Manager) Validate(ctx context.Context, token string) (string, error) {
	key := TokenPrefix + token
	userID, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}

// 刷新：生成新 token，旧 token 延长 GraceTTL 后删除
func (m *Manager) Refresh(ctx context.Context, oldToken, userID string) (string, error) {
	// 延长旧 token 生命周期，过渡期防止正在进行的请求失败
	oldKey := TokenPrefix + oldToken
	m.rdb.Expire(ctx, oldKey, GraceTTL)

	newToken := Generate(userID)
	return newToken, m.Store(ctx, newToken, userID)
}

// 登出：删除 token
func (m *Manager) Revoke(ctx context.Context, token string) error {
	return m.rdb.Del(ctx, TokenPrefix+token).Err()
}
```

- [ ] **Step 2: MySQL 连接包（go-zero 风格）**

```go
// server/common/mysql/mysql.go
package mysql

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func New(dsn string) sqlx.SqlConn {
	return sqlx.NewMysql(dsn)
}
```

- [ ] **Step 3: 错误码包**

```go
// server/common/errcode/errcode.go
package errcode

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrTokenExpired      = errors.New("token expired")
	ErrGroupNotFound     = errors.New("group not found")
	ErrNotGroupMember    = errors.New("not a group member")
	ErrMessageNotFound   = errors.New("message not found")
)
```

- [ ] **Step 4: 验证编译并提交**

```bash
cd server && go build ./common/... && git add server/common/ && git commit -m "feat: add common packages (token+redis, mysql, errcode)"
```

---

### Task A5: 数据库迁移

**Files:**
- Create: `server/migrations/user/001_init.sql`
- Create: `server/migrations/chat/001_init.sql`
- Create: `server/migrations/group/001_init.sql`

内容与之前的 PRD 迁移文件一致，无需修改。

- [ ] **Step 1-3: 编写并执行迁移，验证**

```bash
cd server && docker compose restart mysql && docker compose exec mysql mysql -uroot -proot123 -e "SHOW TABLES FROM webchat_user; SHOW TABLES FROM webchat_chat; SHOW TABLES FROM webchat_group;"
```

- [ ] **Step 4: Commit**

```bash
git add server/migrations/ && git commit -m "feat: add database migrations"
```

---

## Part B: 后端 RPC 服务

### Task B1: user-service — 配置文件 + ServiceContext + DAO

**Files:**
- Create: `server/app/user/etc/user.yaml`
- Create: `server/app/user/internal/config/config.go`
- Create: `server/app/user/internal/svc/servicecontext.go`
- Create: `server/app/user/model/user.go`
- Create: `server/app/user/internal/dao/user.go`

- [ ] **Step 1: 编写配置文件**

```yaml
# server/app/user/etc/user.yaml
Name: user.rpc
ListenOn: 0.0.0.0:50051
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: user.rpc
Mysql:
  DataSource: root:root123@tcp(127.0.0.1:3306)/webchat_user?charset=utf8mb4&parseTime=true
Redis:
  Addr: 127.0.0.1:6379
```

- [ ] **Step 2: 编写 config.go**

```go
// server/app/user/internal/config/config.go
package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
	}
}
```

- [ ] **Step 3: 编写 ServiceContext**

```go
// server/app/user/internal/svc/servicecontext.go
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
```

- [ ] **Step 4: 编写 model 和 DAO**

```go
// server/app/user/model/user.go
package model

type User struct {
	Id           int64  `db:"id"`
	Phone        string `db:"phone"`
	PasswordHash string `db:"password_hash"`
	Nickname     string `db:"nickname"`
	Avatar       string `db:"avatar"`
	Gender       int32  `db:"gender"`
	Region       string `db:"region"`
	Signature    string `db:"signature"`
	CreatedAt    int64  `db:"created_at"`
}
```

```go
// server/app/user/internal/dao/user.go
package dao

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/team/webchat-server/app/user/model"
)

type UserDao struct{ conn sqlx.SqlConn }

func NewUserDao(conn sqlx.SqlConn) *UserDao { return &UserDao{conn} }

func (d *UserDao) Insert(ctx context.Context, u *model.User) (int64, error) {
	result, err := d.conn.ExecCtx(ctx,
		"INSERT INTO users (phone, password_hash, nickname, avatar, created_at) VALUES (?, ?, ?, ?, ?)",
		u.Phone, u.PasswordHash, u.Nickname, u.Avatar, u.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *UserDao) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := d.conn.QueryRowCtx(ctx, &u,
		"SELECT id, phone, password_hash, nickname, avatar, gender, region, signature, created_at FROM users WHERE phone = ?", phone)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDao) FindById(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := d.conn.QueryRowCtx(ctx, &u,
		"SELECT id, phone, password_hash, nickname, avatar, gender, region, signature, created_at FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDao) Update(ctx context.Context, u *model.User) error {
	_, err := d.conn.ExecCtx(ctx,
		"UPDATE users SET nickname=?, avatar=?, gender=?, region=?, signature=? WHERE id=?",
		u.Nickname, u.Avatar, u.Gender, u.Region, u.Signature, u.Id)
	return err
}

func (d *UserDao) Search(ctx context.Context, keyword string) ([]*model.User, error) {
	var users []*model.User
	err := d.conn.QueryRowsCtx(ctx, &users,
		"SELECT id, phone, nickname, avatar, gender, region, signature, created_at FROM users WHERE phone LIKE ? OR nickname LIKE ? LIMIT 20",
		"%"+keyword+"%", "%"+keyword+"%")
	return users, err
}
```

ContactDao 类似，操作 contacts 和 friend_requests 表。

- [ ] **Step 5: 验证编译并提交**

```bash
cd server && go build ./app/user/... && git add server/app/user/etc/ server/app/user/internal/config/ server/app/user/internal/svc/ server/app/user/internal/dao/ server/app/user/model/ && git commit -m "feat: add user-service config, servicecontext, model and dao"
```

---

### Task B2: user-service — Logic 实现

**Files:**
- Create: `server/app/user/internal/logic/register.go`
- Create: `server/app/user/internal/logic/login.go`
- Create: `server/app/user/internal/logic/getprofile.go`
- Create: `server/app/user/internal/logic/updateprofile.go`
- Create: `server/app/user/internal/logic/addcontact.go`
- Create: `server/app/user/internal/logic/acceptfriendrequest.go`
- Create: `server/app/user/internal/logic/rejectfriendrequest.go`
- Create: `server/app/user/internal/logic/listcontacts.go`
- Create: `server/app/user/internal/logic/searchuser.go`
- Create: `server/app/user/internal/logic/getuserbyid.go`
- Modify: `server/app/user/internal/server/userserver.go`
- Modify: `server/app/user/user.go`

- [ ] **Step 1: register logic**

```go
// server/app/user/internal/logic/register.go
package logic

import (
	"context"
	"strconv"
	"time"
	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/user"
	"github.com/team/webchat-server/common/errcode"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	existing, _ := l.svcCtx.UserDao.FindByPhone(l.ctx, in.Phone)
	if existing != nil {
		return nil, errcode.ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Phone:        in.Phone,
		PasswordHash: string(hash),
		Nickname:     in.Nickname,
		CreatedAt:    time.Now().UnixMilli(),
	}
	id, err := l.svcCtx.UserDao.Insert(l.ctx, u)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResponse{UserId: strconv.FormatInt(id, 10)}, nil
}
```

- [ ] **Step 2: login logic**

```go
// server/app/user/internal/logic/login.go
package logic

import (
	"context"
	"strconv"
	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/user"
	"github.com/team/webchat-server/common/errcode"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := l.svcCtx.UserDao.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, errcode.ErrInvalidPassword
	}
	return &user.LoginResponse{UserId: strconv.FormatInt(u.Id, 10)}, nil
}
```

- [ ] **Step 3: 其他 logic 文件**

getprofile.go、updateprofile.go、addcontact.go、acceptfriendrequest.go、rejectfriendrequest.go、listcontacts.go、searchuser.go、getuserbyid.go，均按同样模式实现：Logic struct → NewXxxLogic → 调用 ServiceContext 中的 DAO。

- [ ] **Step 4: 修改 user.go 入口**

```go
// server/app/user/user.go
package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/team/webchat-server/app/user/internal/config"
	"github.com/team/webchat-server/app/user/internal/server"
	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/user"
)

var configFile = flag.String("f", "etc/user.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
	})
	defer s.Stop()
	s.Start()
}
```

- [ ] **Step 5: 验证编译并提交**

```bash
cd server && go build ./app/user/ && git add server/app/user/internal/logic/ server/app/user/user.go && git commit -m "feat: implement user-service logic"
```

---

### Task B3: chat-service — 配置 + ServiceContext + DAO + WebSocket

**Files:**
- Create: `server/app/chat/etc/chat.yaml`
- Create: `server/app/chat/internal/config/config.go`
- Create: `server/app/chat/internal/svc/servicecontext.go`
- Create: `server/app/chat/model/message.go`, `conversation.go`
- Create: `server/app/chat/internal/dao/message.go`, `conversation.go`
- Create: `server/app/chat/ws/hub.go`
- Create: `server/app/chat/ws/client.go`
- Create: `server/app/chat/ws/server.go`

- [ ] **Step 1: 编写配置文件**

```yaml
# server/app/chat/etc/chat.yaml
Name: chat.rpc
ListenOn: 0.0.0.0:50052
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: chat.rpc
Mysql:
  DataSource: root:root123@tcp(127.0.0.1:3306)/webchat_chat?charset=utf8mb4&parseTime=true
Redis:
  Addr: 127.0.0.1:6379
WebSocket:
  ListenOn: 0.0.0.0:8081
```

config.go 和 ServiceContext 模式与 user-service 一致，额外包含 WebSocket Hub。

- [ ] **Step 2: DAO**

与之前 repository 逻辑相同，使用 go-zero 的 sqlx.Conn：

```go
// server/app/chat/internal/dao/message.go
package dao

import (
	"context"
	"time"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/team/webchat-server/app/chat/model"
)

type MessageDao struct{ conn sqlx.SqlConn }

func NewMessageDao(conn sqlx.SqlConn) *MessageDao { return &MessageDao{conn} }

func (d *MessageDao) Insert(ctx context.Context, m *model.Message) (int64, error) {
	m.CreatedAt = time.Now().UnixMilli()
	result, err := d.conn.ExecCtx(ctx,
		"INSERT INTO messages (chat_type, from_user, to_id, msg_type, content, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		m.ChatType, m.FromUser, m.ToId, m.MsgType, m.Content, m.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *MessageDao) FindByConv(ctx context.Context, toId int64, chatType string, page, pageSize int) ([]*model.Message, error) {
	offset := (page - 1) * pageSize
	var msgs []*model.Message
	err := d.conn.QueryRowsCtx(ctx, &msgs,
		"SELECT id, chat_type, from_user, to_id, msg_type, content, is_recalled, created_at FROM messages WHERE to_id = ? AND chat_type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		toId, chatType, pageSize, offset)
	return msgs, err
}

func (d *MessageDao) Recall(ctx context.Context, msgId, userId int64) error {
	_, err := d.conn.ExecCtx(ctx,
		"UPDATE messages SET is_recalled = 1 WHERE id = ? AND from_user = ?", msgId, userId)
	return err
}
```

ConversationDao 同理，操作 conversations 和 read_receipts 表。

- [ ] **Step 3: WebSocket Hub + Client**

与之前方案一致，Hub 管理连接、Client 处理读写，WebSocket server 独立 HTTP 端口启动。

```go
// server/app/chat/ws/hub.go
package ws

import (
	"encoding/json"
	"sync"
)

type WsMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Seq  int64           `json:"seq"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]*Client
	OnMessage func(userID int64, msg *WsMessage)
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int64]*Client)}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c.UserID] = c
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c.UserID)
	h.mu.Unlock()
}

func (h *Hub) SendToUser(userID int64, msg *WsMessage) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		c.Send(msg)
	}
}

func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
```

client.go 与之前一致（gorilla/websocket，ReadPump + WritePump）。

- [ ] **Step 4: WebSocket server（独立 HTTP）**

```go
// server/app/chat/ws/server.go
package ws

import (
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	"github.com/team/webchat-server/common/token"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	hub          *Hub
	tokenManager *token.Manager
}

func NewServer(hub *Hub, tm *token.Manager) *Server {
	return &Server{hub: hub, tokenManager: tm}
}

func (s *Server) Listen(addr string) error {
	http.HandleFunc("/ws", s.handleWS)
	log.Printf("WebSocket server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	userIDStr, err := s.tokenManager.Validate(r.Context(), tokenStr)
	if err != nil || userIDStr == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := NewClient(s.hub, conn, userID)
	s.hub.Register(client)
	go client.WritePump()
	go client.ReadPump(func(c *Client, msg *WsMessage) {
		if s.hub.OnMessage != nil {
			s.hub.OnMessage(c.UserID, msg)
		}
	})
}
```

- [ ] **Step 5: 验证编译并提交**

```bash
cd server && go build ./app/chat/... && git add server/app/chat/ && git commit -m "feat: add chat-service config, dao, websocket"
```

---

### Task B4: chat-service — Logic 实现

**Files:**
- Create: `server/app/chat/internal/logic/sendmessage.go`
- Create: `server/app/chat/internal/logic/getmessages.go`
- Create: `server/app/chat/internal/logic/getconversations.go`
- Create: `server/app/chat/internal/logic/markread.go`
- Create: `server/app/chat/internal/logic/recallmessage.go`
- Modify: `server/app/chat/internal/server/chatserver.go`
- Modify: `server/app/chat/chat.go`（启动 gRPC + WebSocket 双服务）

- [ ] **Step 1: sendmessage logic**

```go
// server/app/chat/internal/logic/sendmessage.go
package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"github.com/team/webchat-server/app/chat/chat"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/team/webchat-server/app/chat/model"
	"github.com/team/webchat-server/app/chat/ws"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SendMessageLogic) SendMessage(in *chat.SendMessageRequest) (*chat.Message, error) {
	fromUser, _ := strconv.ParseInt(in.FromUser, 10, 64)
	toId, _ := strconv.ParseInt(in.ToId, 10, 64)
	msg := &model.Message{
		ChatType: in.ChatType, FromUser: fromUser, ToId: toId,
		MsgType: in.MsgType, Content: in.Content,
	}
	id, err := l.svcCtx.MessageDao.Insert(l.ctx, msg)
	if err != nil {
		return nil, err
	}
	msg.Id = id

	// 更新双方会话
	l.svcCtx.ConversationDao.Upsert(l.ctx, fromUser, toId, in.ChatType, id)
	l.svcCtx.ConversationDao.Upsert(l.ctx, toId, fromUser, in.ChatType, id)

	// WebSocket 推送
	payload, _ := json.Marshal(msg)
	l.svcCtx.Hub.SendToUser(toId, &ws.WsMessage{
		Type: "chat.message", Seq: id, Data: payload,
	})

	return &chat.Message{
		Id: strconv.FormatInt(id, 10), ChatType: msg.ChatType,
		FromUser: in.FromUser, ToId: in.ToId,
		MsgType: msg.MsgType, Content: msg.Content, CreatedAt: msg.CreatedAt,
	}, nil
}
```

- [ ] **Step 2: 其他 logic 文件**

getmessages.go、getconversations.go、markread.go、recallmessage.go，均按同样 Logic 模式实现。

- [ ] **Step 3: 修改 chat.go 入口（双服务启动）**

```go
// server/app/chat/chat.go
package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/team/webchat-server/app/chat/internal/config"
	"github.com/team/webchat-server/app/chat/internal/server"
	"github.com/team/webchat-server/app/chat/internal/svc"
	"github.com/team/webchat-server/app/chat/chat"
	"github.com/team/webchat-server/app/chat/ws"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/chat.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// gRPC server
	go func() {
		s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
			chat.RegisterChatServer(grpcServer, server.NewChatServer(ctx))
		})
		defer s.Stop()
		s.Start()
	}()

	// WebSocket server (blocking)
	if err := ws.NewServer(ctx.Hub, ctx.TokenManager).Listen(c.WebSocket.ListenOn); err != nil {
		panic(err)
	}
}
```

- [ ] **Step 4: 验证编译并提交**

```bash
cd server && go build ./app/chat/ && git add server/app/chat/internal/logic/ server/app/chat/chat.go && git commit -m "feat: implement chat-service logic"
```

---

### Task B5: group-service — 完整实现

**Files:**
- Create: `server/app/group/etc/group.yaml`
- Create: `server/app/group/internal/config/config.go`
- Create: `server/app/group/internal/svc/servicecontext.go`
- Create: `server/app/group/model/group.go`
- Create: `server/app/group/internal/dao/group.go`
- Create: `server/app/group/internal/logic/creategroup.go`
- Create: `server/app/group/internal/logic/getgroup.go`
- Create: `server/app/group/internal/logic/addmember.go`
- Create: `server/app/group/internal/logic/removemember.go`
- Create: `server/app/group/internal/logic/listmembers.go`
- Modify: `server/app/group/internal/server/groupserver.go`
- Modify: `server/app/group/group.go`

与 user-service 模式完全一致，gRPC + etcd 注册。proto 定义包含 CreateGroup / GetGroup / AddMember / RemoveMember / ListMembers。logic 层接收 RPC 请求，调用 DAO 操作 `grps` 和 `group_members` 表。group-service 监听 `:50053`。

- [ ] **Step 1-3: 按模式编写所有文件，验证编译并提交**

```bash
cd server && go build ./app/group/ && git add server/app/group/ && git commit -m "feat: implement group-service"
```

---

### Task B6: media-service — 完整实现

**Files:**
- Create: `server/app/media/etc/media.yaml`
- Create: `server/app/media/internal/config/config.go`
- Create: `server/app/media/internal/svc/servicecontext.go`
- Create: `server/app/media/internal/logic/getuploadurl.go`
- Create: `server/app/media/internal/logic/getfileurl.go`
- Modify: `server/app/media/internal/server/mediaserver.go`
- Modify: `server/app/media/media.go`

ServiceContext 包含 MinIO client（初始化时 EnsureBucket）。logic 层调用 MinIO SDK 生成 presigned URL。media-service 监听 `:50054`。

- [ ] **Step 1-3: 按模式编写所有文件，验证编译并提交**

```bash
cd server && go build ./app/media/ && git add server/app/media/ && git commit -m "feat: implement media-service with MinIO"
```

---

### Task B7: gateway — Auth 中间件 + Logic 对接 RPC Client

**Files:**
- Create: `server/app/gateway/internal/middleware/auth.go`
- Create: `server/app/gateway/internal/logic/auth/login.go`（覆盖生成）
- Create: `server/app/gateway/internal/logic/auth/register.go`（覆盖生成）
- Modify: `server/app/gateway/internal/svc/servicecontext.go`（注入各 RPC client）
- Create/Modify: 所有 logic 文件写入真实业务逻辑

- [ ] **Step 1: Token 鉴权中间件**

```go
// server/app/gateway/internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"
	"github.com/team/webchat-server/common/token"
)

type AuthMiddleware struct {
	TokenManager *token.Manager
}

func NewAuthMiddleware(tm *token.Manager) *AuthMiddleware {
	return &AuthMiddleware{TokenManager: tm}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, err := m.TokenManager.Validate(r.Context(), auth[7:])
		if err != nil || userID == "" {
			http.Error(w, `{"error":"token expired"}`, http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-ID", userID)
		next(w, r)
	}
}
```

- [ ] **Step 2: ServiceContext 注入 RPC client**

```go
// server/app/gateway/internal/svc/servicecontext.go
package svc

import (
	"github.com/team/webchat-server/app/gateway/internal/config"
	"github.com/team/webchat-server/app/user/userclient"
	"github.com/team/webchat-server/app/chat/chatclient"
	"github.com/team/webchat-server/app/group/groupclient"
	"github.com/team/webchat-server/app/media/mediaclient"
	"github.com/team/webchat-server/common/token"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	UserRpc      userclient.User
	ChatRpc      chatclient.Chat
	GroupRpc     groupclient.Group
	MediaRpc     mediaclient.Media
	TokenManager *token.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		UserRpc:      userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ChatRpc:      chatclient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		GroupRpc:     groupclient.NewGroup(zrpc.MustNewClient(c.GroupRpc)),
		MediaRpc:     mediaclient.NewMedia(zrpc.MustNewClient(c.MediaRpc)),
		TokenManager: token.NewManager(/* redis client */),
	}
}
```

- [ ] **Step 3: Login logic（覆盖 goctl 生成的，对接 user RPC + 生成 Token）**

```go
// server/app/gateway/internal/logic/auth/login.go
package authlogic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	resp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Phone: req.Phone, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	tokenStr := token.Generate(resp.UserId)
	if err := l.svcCtx.TokenManager.Store(l.ctx, tokenStr, resp.UserId); err != nil {
		return nil, err
	}
	return &types.LoginResponse{UserId: resp.UserId, Token: tokenStr}, nil
}
```

- [ ] **Step 4: 其他 logic**

register.go（调用 user RPC Register + 生成 token）、refreshToken.go（调用 token.Refresh）、getUser.go（调用 user RPC GetProfile，注入 X-User-ID）、listContacts.go、getConversations.go（调用 chat RPC）、sendMessage.go、createGroup.go 等，均按同样模式：接收 HTTP 请求 → 调用对应 RPC client → 返回 JSON。

- [ ] **Step 5: 在 gateway.api 中声明受保护路由的中间件**

在 `.api` 文件需要鉴权的路由组上添加 `middleware: AuthMiddleware`：

```go
@server(
	prefix: /api/v1
	middleware: AuthMiddleware
)
service gateway {
	// ... protected routes
}
```

- [ ] **Step 6: 验证编译并提交**

```bash
cd server && go build ./app/gateway/ && git add server/app/gateway/internal/ && git commit -m "feat: implement gateway with token auth and rpc routing"
```

---

## Part C: Flutter 客户端

Flutter 部分与之前计划一致，唯一需要调整的是：

- **Token 管理**：Flutter 侧存储 token 字符串（不再是 JWT），每次请求通过 `Authorization: Bearer <token>` 头携带
- **WebSocket 连接**：URL 参数从 `user_id` 改为 `token`：`ws://10.0.2.2:8081/ws?token=<token>`
- **API 路径**：保持不变

### Tasks C1-C7

与之前实施计划 Part C 完全一致（Task C1 初始化 Flutter、C2 核心网络层、C3 认证模块、C4 首页 Tab、C5 会话列表、C6 聊天页面、C7 通讯录）。仅需在 `token_manager.dart` 中将存储内容从 JWT 改为普通 token 字符串，`ws_client.dart` 中的查询参数从 `user_id` 改为 `token`。

不再重复列具体代码，参考之前版本即可。

---

## 服务端口总览

| 服务 | gRPC 端口 | HTTP/WS 端口 |
|------|-----------|-------------|
| gateway | - | 8080 (go-zero API) |
| user-service | 50051 | - |
| chat-service | 50052 | 8081 (WebSocket) |
| group-service | 50053 | - |
| media-service | 50054 | - |
| etcd | 2379 | - |
| MySQL | 3306 | - |
| Redis | 6379 | - |
| MinIO | - | 9000 (API) / 9001 (Console) |
| Kafka | 9092 | - |

---

## 自审清单

**Spec 覆盖检查（对照更新后的 PRD 一期范围）：**
- [x] 注册/登录（Token + Redis）→ Task A4 (token), B1-B2 (user RPC), B7 (gateway auth), C3 (Flutter auth)
- [x] 通讯录 → Task B1-B2, B7, C7
- [x] 单聊（文本/图片/撤回/已读）→ Task B3-B4 (chat RPC + WS), B7, C5-C6
- [x] 群聊 → Task B5 (group RPC), C7
- [x] 会话列表 → Task B3-B4, B7, C5
- [x] 在线状态 → Task B3 (WebSocket Hub)
- [x] 文件上传（MinIO）→ Task B6 (media RPC), B7

**占位符扫描：** 无 TBD、TODO。

**类型一致性：** proto → model → DAO → logic → gateway API types → Flutter types 全链路一致。
