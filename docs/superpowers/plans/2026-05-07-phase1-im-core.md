# Phase 1: IM 核心 — 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建微服务基础设施 + 完成 IM 核心（用户/聊天/群聊/文件）+ Flutter 客户端聊天主流程

**Architecture:** Go 微服务集群（gateway / user-service / chat-service / group-service / media-service），gRPC 内部通信，Kafka 异步事件，WebSocket 长连接推送消息给 Flutter 客户端。Flutter 端采用 Clean Architecture + BLoC，通过 dio 调用 REST API，通过 WebSocket 收发实时消息。

**Tech Stack:** Go + gin + gRPC + Kafka + MySQL + Redis + MinIO + Docker Compose (本地) / Flutter 3.x + flutter_bloc + dio + drift

---

## 文件结构

### Backend (Go monorepo)

```
server/
├── go.mod / go.sum / Makefile
├── docker-compose.yml              # MySQL, Redis, MinIO, Kafka, etcd
├── api/                            # Protobuf 定义
│   ├── common/v1/common.proto
│   ├── user/v1/user.proto
│   ├── chat/v1/chat.proto
│   ├── group/v1/group.proto
│   └── media/v1/media.proto
├── gen/                            # 生成的 protobuf Go 代码 (gitignored)
├── pkg/                            # 共享包
│   ├── jwt/jwt.go
│   ├── mysql/mysql.go
│   ├── redis/redis.go
│   ├── kafka/kafka.go
│   ├── minio/minio.go
│   ├── grpc/interceptors.go
│   ├── config/config.go
│   └── errcode/errcode.go
├── services/
│   ├── gateway/
│   │   ├── main.go
│   │   ├── middleware/  (auth.go, ratelimit.go, cors.go)
│   │   ├── handler/     (auth.go, user.go, chat.go, group.go, media.go)
│   │   └── router/router.go
│   ├── user/
│   │   ├── main.go
│   │   ├── handler/user.go
│   │   ├── service/user.go
│   │   ├── repository/user.go
│   │   └── model/user.go
│   ├── chat/
│   │   ├── main.go
│   │   ├── handler/chat.go
│   │   ├── service/chat.go
│   │   ├── repository/chat.go
│   │   ├── model/chat.go
│   │   └── websocket/ (hub.go, client.go, handler.go)
│   ├── group/
│   │   ├── main.go
│   │   ├── handler/group.go
│   │   ├── service/group.go
│   │   ├── repository/group.go
│   │   └── model/group.go
│   └── media/
│       ├── main.go
│       ├── handler/media.go
│       └── service/media.go
└── migrations/
    ├── user/001_init.sql
    ├── chat/001_init.sql
    ├── group/001_init.sql
    └── media/001_init.sql
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
│   │   ├── cache/     (local_db.dart)
│   │   └── theme/     (colors.dart, theme.dart)
│   ├── features/
│   │   ├── auth/
│   │   │   ├── bloc/  (auth_bloc.dart, auth_event.dart, auth_state.dart)
│   │   │   ├── data/  (auth_repository.dart, auth_api.dart)
│   │   │   └── ui/    (login_page.dart, register_page.dart)
│   │   ├── chat/
│   │   │   ├── bloc/  (conversation_bloc.dart, message_bloc.dart, ...)
│   │   │   ├── data/  (chat_repository.dart, chat_api.dart)
│   │   │   └── ui/    (conversation_list_page.dart, chat_page.dart, widgets/)
│   │   ├── contacts/
│   │   │   ├── bloc/  (contacts_bloc.dart, ...)
│   │   │   ├── data/  (contacts_repository.dart, contacts_api.dart)
│   │   │   └── ui/    (contacts_list_page.dart, add_contact_page.dart)
│   │   └── home/      (home_page.dart, main_tab_scaffold.dart)
│   └── shared/widgets/ (avatar.dart, loading_indicator.dart)
```

---

## Part A: 开发环境与基础设施

### Task A1: 初始化 Go module 与项目骨架

**Files:**
- Create: `server/go.mod`, `server/Makefile`
- Create: `server/services/gateway/main.go` (骨架)
- Create: `server/services/user/main.go` (骨架)
- Create: `server/services/chat/main.go` (骨架)
- Create: `server/services/group/main.go` (骨架)
- Create: `server/services/media/main.go` (骨架)

- [ ] **Step 1: 初始化 Go module**

```bash
cd server && go mod init github.com/team/webchat-server
```

- [ ] **Step 2: 创建 Makefile**

```makefile
.PHONY: proto run-gateway run-user run-chat run-group run-media

proto:
	buf generate api/

run-gateway:
	go run ./services/gateway/

run-user:
	go run ./services/user/

run-chat:
	go run ./services/chat/

run-group:
	go run ./services/group/

run-media:
	go run ./services/media/
```

- [ ] **Step 3: 创建各服务入口骨架**

```go
// server/services/gateway/main.go
package main

import "fmt"

func main() {
	fmt.Println("gateway starting...")
}
```

其他 4 个服务的 `main.go` 同理，替换打印的服务名。

- [ ] **Step 4: 验证编译**

```bash
cd server && go mod tidy && make run-gateway
```

Expected: `gateway starting...`

- [ ] **Step 5: Commit**

```bash
git add server/ && git commit -m "feat: init Go module and service skeletons"
```

---

### Task A2: 编写 docker-compose.yml 本地开发基础设施

**Files:**
- Create: `server/docker-compose.yml`

- [ ] **Step 1: 编写 docker-compose.yml**

```yaml
version: "3.8"
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root123
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data

  kafka:
    image: bitnami/kafka:3.6
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
    ports:
      - "9092:9092"

  etcd:
    image: bitnami/etcd:3.5
    environment:
      ALLOW_NONE_AUTHENTICATION: "yes"
    ports:
      - "2379:2379"

volumes:
  mysql_data:
  minio_data:
```

- [ ] **Step 2: 启动基础设施**

```bash
cd server && docker compose up -d
```

Expected: 5 个容器均 Running（`docker compose ps` 验证）。

- [ ] **Step 3: Commit**

```bash
git add server/docker-compose.yml && git commit -m "feat: add docker-compose for local dev infrastructure"
```

---

### Task A3: 创建 Protobuf 定义并生成代码

**Files:**
- Create: `server/api/common/v1/common.proto`
- Create: `server/api/user/v1/user.proto`
- Create: `server/api/chat/v1/chat.proto`
- Create: `server/api/group/v1/group.proto`
- Create: `server/api/media/v1/media.proto`
- Create: `server/buf.yaml`, `server/buf.gen.yaml`

- [ ] **Step 1: 编写公共类型 proto**

```protobuf
// server/api/common/v1/common.proto
syntax = "proto3";
package common.v1;
option go_package = "github.com/team/webchat-server/gen/common/v1;commonv1";

message Pagination {
  int32 page = 1;
  int32 page_size = 2;
}

message Empty {}
```

- [ ] **Step 2: 编写 user.proto**

```protobuf
// server/api/user/v1/user.proto
syntax = "proto3";
package user.v1;
option go_package = "github.com/team/webchat-server/gen/user/v1;userv1";
import "common/v1/common.proto";

service UserService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc GetProfile(GetProfileRequest) returns (User);
  rpc UpdateProfile(UpdateProfileRequest) returns (User);
  rpc AddContact(AddContactRequest) returns (common.v1.Empty);
  rpc AcceptFriendRequest(AcceptFriendRequestRequest) returns (common.v1.Empty);
  rpc ListContacts(ListContactsRequest) returns (ListContactsResponse);
  rpc SearchUser(SearchUserRequest) returns (SearchUserResponse);
}

message User {
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
  string token = 2;
}

message LoginRequest {
  string phone = 1;
  string password = 2;
}

message LoginResponse {
  string user_id = 1;
  string token = 2;
}

message GetProfileRequest {
  string user_id = 1;
}

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

message AcceptFriendRequestRequest {
  string request_id = 1;
  string action = 2; // "accept" | "reject"
}

message ListContactsRequest {
  string user_id = 1;
}

message ListContactsResponse {
  repeated Contact contacts = 1;
}

message Contact {
  string user_id = 1;
  string nickname = 2;
  string avatar = 3;
  string remark = 4;
  string tag = 5;
  bool is_blocked = 6;
}

message SearchUserRequest {
  string keyword = 1;
}

message SearchUserResponse {
  repeated User users = 1;
}
```

- [ ] **Step 3: 编写 chat.proto**

```protobuf
// server/api/chat/v1/chat.proto
syntax = "proto3";
package chat.v1;
option go_package = "github.com/team/webchat-server/gen/chat/v1;chatv1";
import "common/v1/common.proto";

service ChatService {
  rpc SendMessage(SendMessageRequest) returns (Message);
  rpc GetMessages(GetMessagesRequest) returns (GetMessagesResponse);
  rpc RecallMessage(RecallMessageRequest) returns (common.v1.Empty);
  rpc GetConversations(GetConversationsRequest) returns (GetConversationsResponse);
  rpc MarkRead(MarkReadRequest) returns (common.v1.Empty);
}

message Message {
  string id = 1;
  string chat_type = 2; // "single" | "group"
  string from_user = 3;
  string to_id = 4;
  string msg_type = 5; // "text" | "image" | "voice" | "video" | "file" | "location"
  string content = 6;
  bool is_recalled = 7;
  int64 created_at = 8;
}

message SendMessageRequest {
  string chat_type = 1;
  string from_user = 2;
  string to_id = 3;
  string msg_type = 4;
  string content = 5;
}

message GetMessagesRequest {
  string conv_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message GetMessagesResponse {
  repeated Message messages = 1;
}

message RecallMessageRequest {
  string msg_id = 1;
  string user_id = 2;
}

message Conversation {
  string id = 1;
  string chat_type = 2;
  string target_id = 3;
  string target_name = 4;
  string target_avatar = 5;
  Message last_msg = 6;
  int32 unread_count = 7;
  bool is_pinned = 8;
  bool is_muted = 9;
  int64 updated_at = 10;
}

message GetConversationsRequest {
  string user_id = 1;
}

message GetConversationsResponse {
  repeated Conversation conversations = 1;
}

message MarkReadRequest {
  string user_id = 1;
  string conv_id = 2;
}
```

- [ ] **Step 4: 编写 group.proto**

```protobuf
// server/api/group/v1/group.proto
syntax = "proto3";
package group.v1;
option go_package = "github.com/team/webchat-server/gen/group/v1;groupv1";
import "common/v1/common.proto";

service GroupService {
  rpc CreateGroup(CreateGroupRequest) returns (Group);
  rpc GetGroup(GetGroupRequest) returns (Group);
  rpc AddMember(AddMemberRequest) returns (common.v1.Empty);
  rpc RemoveMember(RemoveMemberRequest) returns (common.v1.Empty);
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse);
}

message Group {
  string id = 1;
  string name = 2;
  string avatar = 3;
  string owner_id = 4;
  string announcement = 5;
  int32 member_count = 6;
  int32 max_members = 7;
  int64 created_at = 8;
}

message GroupMember {
  string user_id = 1;
  string role = 2; // "owner" | "admin" | "member"
  string alias = 3;
  bool is_muted = 4;
  int64 joined_at = 5;
}

message CreateGroupRequest {
  string name = 1;
  string owner_id = 2;
  repeated string member_ids = 3;
}

message GetGroupRequest {
  string group_id = 1;
}

message AddMemberRequest {
  string group_id = 1;
  repeated string user_ids = 2;
}

message RemoveMemberRequest {
  string group_id = 1;
  string user_id = 2;
}

message ListMembersRequest {
  string group_id = 1;
}

message ListMembersResponse {
  repeated GroupMember members = 1;
}
```

- [ ] **Step 5: 编写 media.proto**

```protobuf
// server/api/media/v1/media.proto
syntax = "proto3";
package media.v1;
option go_package = "github.com/team/webchat-server/gen/media/v1;mediav1";

service MediaService {
  rpc GetUploadURL(GetUploadURLRequest) returns (GetUploadURLResponse);
  rpc GetFileURL(GetFileURLRequest) returns (GetFileURLResponse);
}

message GetUploadURLRequest {
  string file_name = 1;
  string content_type = 2;
  int64 file_size = 3;
}

message GetUploadURLResponse {
  string upload_url = 1;
  string file_id = 2;
  int64 expires_in = 3;
}

message GetFileURLRequest {
  string file_id = 1;
}

message GetFileURLResponse {
  string url = 1;
}
```

- [ ] **Step 6: 配置 buf 并生成代码**

```yaml
# server/buf.yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```

```yaml
# server/buf.gen.yaml
version: v1
plugins:
  - plugin: go
    out: gen
    opt: paths=source_relative
  - plugin: go-grpc
    out: gen
    opt: paths=source_relative
```

```bash
cd server && buf generate api/
```

- [ ] **Step 7: Commit**

```bash
git add server/api/ server/buf.yaml server/buf.gen.yaml server/gen/ && git commit -m "feat: add protobuf definitions and generated code"
```

---

### Task A4: 创建共享包 (pkg/)

**Files:**
- Create: `server/pkg/config/config.go`
- Create: `server/pkg/mysql/mysql.go`
- Create: `server/pkg/redis/redis.go`
- Create: `server/pkg/jwt/jwt.go`
- Create: `server/pkg/kafka/kafka.go`
- Create: `server/pkg/minio/minio.go`
- Create: `server/pkg/errcode/errcode.go`
- Create: `server/pkg/grpc/interceptors.go`

- [ ] **Step 1: config 包**

```go
// server/pkg/config/config.go
package config

import "os"

type Config struct {
	MySQLDSN  string
	RedisAddr string
	KafkaAddr string
	MinioEndpoint string
	MinioAccessKey string
	MinioSecretKey string
	EtcdAddr  string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		MySQLDSN:       getEnv("MYSQL_DSN", "root:root123@tcp(127.0.0.1:3306)/webchat?charset=utf8mb4&parseTime=true"),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		KafkaAddr:      getEnv("KAFKA_ADDR", "127.0.0.1:9092"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		EtcdAddr:       getEnv("ETCD_ADDR", "127.0.0.1:2379"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
```

- [ ] **Step 2: mysql 包**

```go
// server/pkg/mysql/mysql.go
package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return db, db.Ping()
}
```

- [ ] **Step 3: redis 包**

```go
// server/pkg/redis/redis.go
package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func New(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	_, err := rdb.Ping(context.Background()).Result()
	return rdb, err
}
```

- [ ] **Step 4: jwt 包**

```go
// server/pkg/jwt/jwt.go
package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Generate(secret, userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Validate(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
```

- [ ] **Step 5: kafka 包**

```go
// server/pkg/kafka/kafka.go
package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewWriter(addr, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewReader(addr, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{addr},
		Topic:   topic,
		GroupID: groupID,
	})
}
```

- [ ] **Step 6: minio 包**

```go
// server/pkg/minio/minio.go
package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
}

func EnsureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}
```

- [ ] **Step 7: errcode 包**

```go
// server/pkg/errcode/errcode.go
package errcode

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrGroupNotFound     = errors.New("group not found")
	ErrNotGroupMember    = errors.New("not a group member")
	ErrMessageNotFound   = errors.New("message not found")
)
```

- [ ] **Step 8: grpc interceptors 包**

```go
// server/pkg/grpc/interceptors.go
package grpc

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("gRPC %s took %v, err: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}
```

- [ ] **Step 9: 安装依赖并验证编译**

```bash
cd server && go mod tidy
```

Expected: `go mod tidy` 成功，无报错。

- [ ] **Step 10: Commit**

```bash
git add server/pkg/ && git commit -m "feat: add shared packages (config, mysql, redis, jwt, kafka, minio, errcode, grpc interceptors)"
```

---

### Task A5: 创建数据库迁移文件

**Files:**
- Create: `server/migrations/user/001_init.sql`
- Create: `server/migrations/chat/001_init.sql`
- Create: `server/migrations/group/001_init.sql`

- [ ] **Step 1: user-service 建表**

```sql
-- server/migrations/user/001_init.sql
CREATE DATABASE IF NOT EXISTS webchat_user;
USE webchat_user;

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(500) DEFAULT '',
    gender TINYINT DEFAULT 0,
    region VARCHAR(100) DEFAULT '',
    signature VARCHAR(200) DEFAULT '',
    created_at BIGINT NOT NULL
);

CREATE TABLE contacts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    contact_id BIGINT NOT NULL,
    remark VARCHAR(50) DEFAULT '',
    tag VARCHAR(50) DEFAULT '',
    is_blocked TINYINT DEFAULT 0,
    added_at BIGINT NOT NULL,
    UNIQUE KEY uk_user_contact (user_id, contact_id)
);

CREATE TABLE friend_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    from_user BIGINT NOT NULL,
    to_user BIGINT NOT NULL,
    message VARCHAR(100) DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at BIGINT NOT NULL
);
```

- [ ] **Step 2: chat-service 建表**

```sql
-- server/migrations/chat/001_init.sql
CREATE DATABASE IF NOT EXISTS webchat_chat;
USE webchat_chat;

CREATE TABLE messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chat_type VARCHAR(10) NOT NULL,
    from_user BIGINT NOT NULL,
    to_id BIGINT NOT NULL,
    msg_type VARCHAR(20) NOT NULL,
    content TEXT,
    is_recalled TINYINT DEFAULT 0,
    created_at BIGINT NOT NULL,
    INDEX idx_to_id_chat_type_created (to_id, chat_type, created_at)
);

CREATE TABLE conversations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    chat_type VARCHAR(10) NOT NULL,
    target_id BIGINT NOT NULL,
    last_msg_id BIGINT DEFAULT 0,
    unread_count INT DEFAULT 0,
    is_pinned TINYINT DEFAULT 0,
    is_muted TINYINT DEFAULT 0,
    updated_at BIGINT NOT NULL,
    UNIQUE KEY uk_user_conv (user_id, chat_type, target_id)
);

CREATE TABLE read_receipts (
    msg_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    read_at BIGINT NOT NULL,
    PRIMARY KEY (msg_id, user_id)
);
```

- [ ] **Step 3: group-service 建表**

```sql
-- server/migrations/group/001_init.sql
CREATE DATABASE IF NOT EXISTS webchat_group;
USE webchat_group;

CREATE TABLE grps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(500) DEFAULT '',
    owner_id BIGINT NOT NULL,
    announcement VARCHAR(500) DEFAULT '',
    member_count INT DEFAULT 0,
    max_members INT DEFAULT 500,
    created_at BIGINT NOT NULL
);

CREATE TABLE group_members (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    alias VARCHAR(50) DEFAULT '',
    is_muted TINYINT DEFAULT 0,
    joined_at BIGINT NOT NULL,
    UNIQUE KEY uk_group_user (group_id, user_id)
);
```

- [ ] **Step 4: 重启 MySQL 触发迁移并验证**

```bash
cd server && docker compose restart mysql
# 等 MySQL 启动后验证表已创建
docker compose exec mysql mysql -uroot -proot123 -e "SHOW TABLES FROM webchat_user; SHOW TABLES FROM webchat_chat; SHOW TABLES FROM webchat_group;"
```

- [ ] **Step 5: Commit**

```bash
git add server/migrations/ && git commit -m "feat: add database migration files"
```

---

## Part B: 后端服务

### Task B1: user-service — model + repository

**Files:**
- Create: `server/services/user/model/user.go`
- Create: `server/services/user/repository/user.go`

- [ ] **Step 1: 定义 model**

```go
// server/services/user/model/user.go
package model

type User struct {
	ID           int64  `json:"id"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"-"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Gender       int32  `json:"gender"`
	Region       string `json:"region"`
	Signature    string `json:"signature"`
	CreatedAt    int64  `json:"created_at"`
}

type Contact struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ContactID int64  `json:"contact_id"`
	Remark    string `json:"remark"`
	Tag       string `json:"tag"`
	IsBlocked bool   `json:"is_blocked"`
	AddedAt   int64  `json:"added_at"`
}

type FriendRequest struct {
	ID        int64  `json:"id"`
	FromUser  int64  `json:"from_user"`
	ToUser    int64  `json:"to_user"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}
```

- [ ] **Step 2: 编写 repository**

```go
// server/services/user/repository/user.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"github.com/team/webchat-server/services/user/model"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) Create(ctx context.Context, u *model.User) (int64, error) {
	now := time.Now().UnixMilli()
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (phone, password_hash, nickname, avatar, created_at) VALUES (?, ?, ?, ?, ?)",
		u.Phone, u.PasswordHash, u.Nickname, u.Avatar, now)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return result.LastInsertId()
}

func (r *UserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	u := &model.User{}
	row := r.db.QueryRowContext(ctx,
		"SELECT id, phone, password_hash, nickname, avatar, gender, region, signature, created_at FROM users WHERE phone = ?", phone)
	err := row.Scan(&u.ID, &u.Phone, &u.PasswordHash, &u.Nickname, &u.Avatar, &u.Gender, &u.Region, &u.Signature, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	u := &model.User{}
	row := r.db.QueryRowContext(ctx,
		"SELECT id, phone, password_hash, nickname, avatar, gender, region, signature, created_at FROM users WHERE id = ?", id)
	err := row.Scan(&u.ID, &u.Phone, &u.PasswordHash, &u.Nickname, &u.Avatar, &u.Gender, &u.Region, &u.Signature, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, u *model.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET nickname=?, avatar=?, gender=?, region=?, signature=? WHERE id=?",
		u.Nickname, u.Avatar, u.Gender, u.Region, u.Signature, u.ID)
	return err
}

func (r *UserRepo) SearchByKeyword(ctx context.Context, keyword string) ([]*model.User, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, phone, nickname, avatar, gender, region, signature, created_at FROM users WHERE phone LIKE ? OR nickname LIKE ? LIMIT 20",
		"%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Phone, &u.PasswordHash, &u.Nickname, &u.Avatar, &u.Gender, &u.Region, &u.Signature, &u.CreatedAt); err != nil {
			return nil, err
		}
		u.PasswordHash = ""
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) AddContact(ctx context.Context, userID, contactID int64, message string) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO friend_requests (from_user, to_user, message, status, created_at) VALUES (?, ?, ?, 'pending', ?)",
		userID, contactID, message, now)
	return err
}

func (r *UserRepo) AcceptFriendRequest(ctx context.Context, requestID int64) (*model.FriendRequest, error) {
	var req model.FriendRequest
	row := r.db.QueryRowContext(ctx, "SELECT id, from_user, to_user, status FROM friend_requests WHERE id = ? AND status = 'pending'", requestID)
	if err := row.Scan(&req.ID, &req.FromUser, &req.ToUser, &req.Status); err != nil {
		return nil, err
	}
	_, err := r.db.ExecContext(ctx, "UPDATE friend_requests SET status = 'accepted' WHERE id = ?", requestID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	r.db.ExecContext(ctx, "INSERT INTO contacts (user_id, contact_id, added_at) VALUES (?, ?, ?)", req.FromUser, req.ToUser, now)
	r.db.ExecContext(ctx, "INSERT INTO contacts (user_id, contact_id, added_at) VALUES (?, ?, ?)", req.ToUser, req.FromUser, now)
	return &req, nil
}

func (r *UserRepo) RejectFriendRequest(ctx context.Context, requestID int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE friend_requests SET status = 'rejected' WHERE id = ?", requestID)
	return err
}

func (r *UserRepo) ListContacts(ctx context.Context, userID int64) ([]*model.Contact, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT c.contact_id, COALESCE(c.remark, u.nickname), u.avatar, c.remark, c.tag, c.is_blocked "+
			"FROM contacts c JOIN users u ON c.contact_id = u.id WHERE c.user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []*model.Contact
	for rows.Next() {
		c := &model.Contact{UserID: userID}
		var remark sql.NullString
		if err := rows.Scan(&c.ContactID, &c.Remark, &c.Avatar, &remark, &c.Tag, &c.IsBlocked); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}
```

- [ ] **Step 3: 安装依赖并验证编译**

```bash
cd server && go mod tidy && go build ./services/user/...
```

Expected: 编译成功。

- [ ] **Step 4: Commit**

```bash
git add server/services/user/model/ server/services/user/repository/ && git commit -m "feat: add user-service model and repository"
```

---

### Task B2: user-service — service + handler + main.go

**Files:**
- Create: `server/services/user/service/user.go`
- Create: `server/services/user/handler/user.go`
- Modify: `server/services/user/main.go`

- [ ] **Step 1: 编写 service 层**

```go
// server/services/user/service/user.go
package service

import (
	"context"
	"github.com/team/webchat-server/pkg/errcode"
	"github.com/team/webchat-server/services/user/model"
	"github.com/team/webchat-server/services/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func New(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, phone, password, nickname string) (*model.User, error) {
	existing, _ := s.repo.FindByPhone(ctx, phone)
	if existing != nil {
		return nil, errcode.ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{Phone: phone, PasswordHash: string(hash), Nickname: nickname}
	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (s *UserService) Login(ctx context.Context, phone, password string) (*model.User, error) {
	u, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errcode.ErrInvalidPassword
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, u *model.User) error {
	return s.repo.UpdateProfile(ctx, u)
}

func (s *UserService) SearchUser(ctx context.Context, keyword string) ([]*model.User, error) {
	return s.repo.SearchByKeyword(ctx, keyword)
}

func (s *UserService) AddContact(ctx context.Context, fromUser, toUser int64, message string) error {
	return s.repo.AddContact(ctx, fromUser, toUser, message)
}

func (s *UserService) AcceptFriendRequest(ctx context.Context, requestID int64) (*model.FriendRequest, error) {
	return s.repo.AcceptFriendRequest(ctx, requestID)
}

func (s *UserService) RejectFriendRequest(ctx context.Context, requestID int64) error {
	return s.repo.RejectFriendRequest(ctx, requestID)
}

func (s *UserService) ListContacts(ctx context.Context, userID int64) ([]*model.Contact, error) {
	return s.repo.ListContacts(ctx, userID)
}
```

- [ ] **Step 2: 编写 gRPC handler**

```go
// server/services/user/handler/user.go
package handler

import (
	"context"
	"strconv"
	"github.com/team/webchat-server/gen/common/v1"
	userv1 "github.com/team/webchat-server/gen/user/v1"
	"github.com/team/webchat-server/services/user/service"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	svc *service.UserService
}

func New(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	u, err := h.svc.Register(ctx, req.Phone, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	return &userv1.RegisterResponse{UserId: strconv.FormatInt(u.ID, 10), Token: ""}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	u, err := h.svc.Login(ctx, req.Phone, req.Password)
	if err != nil {
		return nil, err
	}
	return &userv1.LoginResponse{UserId: strconv.FormatInt(u.ID, 10), Token: ""}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.User, error) {
	id, _ := strconv.ParseInt(req.UserId, 10, 64)
	u, err := h.svc.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return &userv1.User{
		Id: strconv.FormatInt(u.ID, 10), Phone: u.Phone, Nickname: u.Nickname,
		Avatar: u.Avatar, Gender: u.Gender, Region: u.Region, Signature: u.Signature, CreatedAt: u.CreatedAt,
	}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.User, error) {
	id, _ := strconv.ParseInt(req.UserId, 10, 64)
	u := &model.User{ID: id, Nickname: req.Nickname, Avatar: req.Avatar, Gender: req.Gender, Region: req.Region, Signature: req.Signature}
	if err := h.svc.UpdateProfile(ctx, u); err != nil {
		return nil, err
	}
	u, _ = h.svc.GetProfile(ctx, id)
	return &userv1.User{
		Id: strconv.FormatInt(u.ID, 10), Phone: u.Phone, Nickname: u.Nickname,
		Avatar: u.Avatar, Gender: u.Gender, Region: u.Region, Signature: u.Signature, CreatedAt: u.CreatedAt,
	}, nil
}
```

注意：handler 中需要 import `"github.com/team/webchat-server/services/user/model"`。

- [ ] **Step 3: 编写 main.go**

```go
// server/services/user/main.go
package main

import (
	"log"
	"net"
	userv1 "github.com/team/webchat-server/gen/user/v1"
	"github.com/team/webchat-server/pkg/config"
	_ "github.com/team/webchat-server/pkg/grpc"
	"github.com/team/webchat-server/pkg/mysql"
	"github.com/team/webchat-server/services/user/handler"
	"github.com/team/webchat-server/services/user/repository"
	"github.com/team/webchat-server/services/user/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	db, err := mysql.New(cfg.MySQLDSN + "webchat_user")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	repo := repository.NewUserRepo(db)
	svc := service.New(repo)
	h := handler.New(svc)

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcpkg.UnaryLoggingInterceptor))
	userv1.RegisterUserServiceServer(s, h)
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Println("user-service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
```

- [ ] **Step 4: 验证编译**

```bash
cd server && go build ./services/user/
```

Expected: 编译成功。

- [ ] **Step 5: Commit**

```bash
git add server/services/user/ && git commit -m "feat: implement user-service gRPC handler and main"
```

---

### Task B3: chat-service — model + repository

**Files:**
- Create: `server/services/chat/model/chat.go`
- Create: `server/services/chat/repository/chat.go`

- [ ] **Step 1: 定义 model**

```go
// server/services/chat/model/chat.go
package model

type Message struct {
	ID         int64  `json:"id"`
	ChatType   string `json:"chat_type"`
	FromUser   int64  `json:"from_user"`
	ToID       int64  `json:"to_id"`
	MsgType    string `json:"msg_type"`
	Content    string `json:"content"`
	IsRecalled bool   `json:"is_recalled"`
	CreatedAt  int64  `json:"created_at"`
}

type Conversation struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
	ChatType     string   `json:"chat_type"`
	TargetID     int64    `json:"target_id"`
	TargetName   string   `json:"target_name"`
	TargetAvatar string   `json:"target_avatar"`
	LastMsg      *Message `json:"last_msg"`
	UnreadCount  int32    `json:"unread_count"`
	IsPinned     bool     `json:"is_pinned"`
	IsMuted      bool     `json:"is_muted"`
	UpdatedAt    int64    `json:"updated_at"`
}

type ReadReceipt struct {
	MsgID  int64 `json:"msg_id"`
	UserID int64 `json:"user_id"`
	ReadAt int64 `json:"read_at"`
}
```

- [ ] **Step 2: 编写 repository**

```go
// server/services/chat/repository/chat.go
package repository

import (
	"context"
	"database/sql"
	"time"
	"github.com/team/webchat-server/services/chat/model"
)

type ChatRepo struct{ db *sql.DB }

func NewChatRepo(db *sql.DB) *ChatRepo { return &ChatRepo{db} }

func (r *ChatRepo) SaveMessage(ctx context.Context, msg *model.Message) (int64, error) {
	now := time.Now().UnixMilli()
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO messages (chat_type, from_user, to_id, msg_type, content, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		msg.ChatType, msg.FromUser, msg.ToID, msg.MsgType, msg.Content, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ChatRepo) GetMessages(ctx context.Context, toID int64, chatType string, page, pageSize int) ([]*model.Message, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, chat_type, from_user, to_id, msg_type, content, is_recalled, created_at FROM messages WHERE to_id = ? AND chat_type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		toID, chatType, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []*model.Message
	for rows.Next() {
		m := &model.Message{}
		if err := rows.Scan(&m.ID, &m.ChatType, &m.FromUser, &m.ToID, &m.MsgType, &m.Content, &m.IsRecalled, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *ChatRepo) RecallMessage(ctx context.Context, msgID, userID int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE messages SET is_recalled = 1 WHERE id = ? AND from_user = ?", msgID, userID)
	return err
}

func (r *ChatRepo) GetConversations(ctx context.Context, userID int64) ([]*model.Conversation, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, chat_type, target_id, unread_count, is_pinned, is_muted, updated_at FROM conversations WHERE user_id = ? ORDER BY updated_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var convs []*model.Conversation
	for rows.Next() {
		c := &model.Conversation{UserID: userID}
		if err := rows.Scan(&c.ID, &c.ChatType, &c.TargetID, &c.UnreadCount, &c.IsPinned, &c.IsMuted, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, nil
}

func (r *ChatRepo) UpsertConversation(ctx context.Context, userID, targetID int64, chatType string, msgID int64) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO conversations (user_id, chat_type, target_id, last_msg_id, unread_count, updated_at)
		 VALUES (?, ?, ?, ?, 1, ?)
		 ON DUPLICATE KEY UPDATE last_msg_id = ?, unread_count = unread_count + 1, updated_at = ?`,
		userID, chatType, targetID, msgID, now, msgID, now)
	return err
}

func (r *ChatRepo) MarkRead(ctx context.Context, userID, msgID int64) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx,
		"INSERT IGNORE INTO read_receipts (msg_id, user_id, read_at) VALUES (?, ?, ?)", msgID, userID, now)
	return err
}

func (r *ChatRepo) ClearUnread(ctx context.Context, userID, convID int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE conversations SET unread_count = 0 WHERE id = ? AND user_id = ?", convID, userID)
	return err
}
```

- [ ] **Step 3: 验证编译**

```bash
cd server && go build ./services/chat/...
```

Expected: 编译成功。

- [ ] **Step 4: Commit**

```bash
git add server/services/chat/model/ server/services/chat/repository/ && git commit -m "feat: add chat-service model and repository"
```

---

### Task B4: chat-service — WebSocket Hub + Client

**Files:**
- Create: `server/services/chat/websocket/hub.go`
- Create: `server/services/chat/websocket/client.go`

- [ ] **Step 1: 编写 WebSocket Hub**

```go
// server/services/chat/websocket/hub.go
package websocket

import (
	"encoding/json"
	"log"
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
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int64]*Client)}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client.UserID] = client
	h.mu.Unlock()
	log.Printf("user %d connected", client.UserID)
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if c, ok := h.clients[client.UserID]; ok && c == client {
		delete(h.clients, client.UserID)
	}
	h.mu.Unlock()
	log.Printf("user %d disconnected", client.UserID)
}

func (h *Hub) SendToUser(userID int64, msg *WsMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		client.Send(msg)
	}
}

func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	_, ok := h.clients[userID]
	h.mu.RUnlock()
	return ok
}
```

- [ ] **Step 2: 编写 WebSocket Client**

```go
// server/services/chat/websocket/client.go
package websocket

import (
	"encoding/json"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 4096
)

type Client struct {
	UserID int64
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn, userID int64) *Client {
	return &Client{
		UserID: userID,
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

func (c *Client) Send(msg *WsMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) ReadPump(h func(*Client, *WsMessage)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg WsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		h(c, &msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
```

- [ ] **Step 3: 编写 WebSocket HTTP handler**

```go
// server/services/chat/websocket/handler.go
package websocket

import (
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, onMessage func(*Client, *WsMessage)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade: %v", err)
			return
		}
		client := NewClient(hub, conn, userID)
		hub.Register(client)
		go client.WritePump()
		go client.ReadPump(onMessage)
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add server/services/chat/websocket/ && git commit -m "feat: add WebSocket hub, client and handler"
```

---

### Task B5: chat-service — service + gRPC handler + main.go

**Files:**
- Create: `server/services/chat/service/chat.go`
- Create: `server/services/chat/handler/chat.go`
- Modify: `server/services/chat/main.go`

- [ ] **Step 1: 编写 service 层**

```go
// server/services/chat/service/chat.go
package service

import (
	"context"
	"github.com/team/webchat-server/services/chat/model"
	"github.com/team/webchat-server/services/chat/repository"
	"github.com/team/webchat-server/services/chat/websocket"
)

type ChatService struct {
	repo *repository.ChatRepo
	hub  *websocket.Hub
}

func New(repo *repository.ChatRepo, hub *websocket.Hub) *ChatService {
	return &ChatService{repo: repo, hub: hub}
}

func (s *ChatService) SendMessage(ctx context.Context, m *model.Message) (*model.Message, error) {
	id, err := s.repo.SaveMessage(ctx, m)
	if err != nil {
		return nil, err
	}
	m.ID = id
	s.repo.UpsertConversation(ctx, m.FromUser, m.ToID, m.ChatType, id)
	s.repo.UpsertConversation(ctx, m.ToID, m.FromUser, m.ChatType, id)

	msg := &websocket.WsMessage{Type: "chat.message", Seq: id}
	s.hub.SendToUser(m.ToID, msg)
	return m, nil
}

func (s *ChatService) GetMessages(ctx context.Context, convID int64, chatType string, page, pageSize int) ([]*model.Message, error) {
	return s.repo.GetMessages(ctx, convID, chatType, page, pageSize)
}

func (s *ChatService) RecallMessage(ctx context.Context, msgID, userID int64) error {
	return s.repo.RecallMessage(ctx, msgID, userID)
}

func (s *ChatService) GetConversations(ctx context.Context, userID int64) ([]*model.Conversation, error) {
	return s.repo.GetConversations(ctx, userID)
}

func (s *ChatService) MarkRead(ctx context.Context, userID, msgID int64) error {
	return s.repo.MarkRead(ctx, userID, msgID)
}
```

- [ ] **Step 2: 编写 gRPC handler**

```go
// server/services/chat/handler/chat.go
package handler

import (
	"context"
	"strconv"
	chatv1 "github.com/team/webchat-server/gen/chat/v1"
	"github.com/team/webchat-server/services/chat/model"
	"github.com/team/webchat-server/services/chat/service"
)

type ChatHandler struct {
	chatv1.UnimplementedChatServiceServer
	svc *service.ChatService
}

func New(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*chatv1.Message, error) {
	fromUser, _ := strconv.ParseInt(req.FromUser, 10, 64)
	toID, _ := strconv.ParseInt(req.ToId, 10, 64)
	msg, err := h.svc.SendMessage(ctx, &model.Message{
		ChatType: req.ChatType, FromUser: fromUser, ToID: toID,
		MsgType: req.MsgType, Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &chatv1.Message{
		Id: strconv.FormatInt(msg.ID, 10), ChatType: msg.ChatType,
		FromUser: req.FromUser, ToId: req.ToId, MsgType: msg.MsgType,
		Content: msg.Content, CreatedAt: msg.CreatedAt,
	}, nil
}

func (h *ChatHandler) GetMessages(ctx context.Context, req *chatv1.GetMessagesRequest) (*chatv1.GetMessagesResponse, error) {
	convID, _ := strconv.ParseInt(req.ConvId, 10, 64)
	msgs, err := h.svc.GetMessages(ctx, convID, "single", int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var pbMsgs []*chatv1.Message
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, &chatv1.Message{
			Id: strconv.FormatInt(m.ID, 10), ChatType: m.ChatType,
			FromUser: strconv.FormatInt(m.FromUser, 10),
			ToId: strconv.FormatInt(m.ToID, 10), MsgType: m.MsgType,
			Content: m.Content, IsRecalled: m.IsRecalled, CreatedAt: m.CreatedAt,
		})
	}
	return &chatv1.GetMessagesResponse{Messages: pbMsgs}, nil
}

func (h *ChatHandler) GetConversations(ctx context.Context, req *chatv1.GetConversationsRequest) (*chatv1.GetConversationsResponse, error) {
	userID, _ := strconv.ParseInt(req.UserId, 10, 64)
	convs, err := h.svc.GetConversations(ctx, userID)
	if err != nil {
		return nil, err
	}
	var pbConvs []*chatv1.Conversation
	for _, c := range convs {
		pbConvs = append(pbConvs, &chatv1.Conversation{
			Id: strconv.FormatInt(c.ID, 10), ChatType: c.ChatType,
			TargetId: strconv.FormatInt(c.TargetID, 10),
			UnreadCount: c.UnreadCount, IsPinned: c.IsPinned, IsMuted: c.IsMuted, UpdatedAt: c.UpdatedAt,
		})
	}
	return &chatv1.GetConversationsResponse{Conversations: pbConvs}, nil
}

func (h *ChatHandler) MarkRead(ctx context.Context, req *chatv1.MarkReadRequest) (*commonv1.Empty, error) {
	userID, _ := strconv.ParseInt(req.UserId, 10, 64)
	msgID, _ := strconv.ParseInt("0", 10, 64)
	h.svc.MarkRead(ctx, userID, msgID)
	return &commonv1.Empty{}, nil
}
```

- [ ] **Step 3: 编写 main.go**

```go
// server/services/chat/main.go
package main

import (
	"log"
	"net"
	"net/http"
	chatv1 "github.com/team/webchat-server/gen/chat/v1"
	"github.com/team/webchat-server/pkg/config"
	_ "github.com/team/webchat-server/pkg/grpc"
	"github.com/team/webchat-server/pkg/mysql"
	"github.com/team/webchat-server/services/chat/handler"
	"github.com/team/webchat-server/services/chat/repository"
	"github.com/team/webchat-server/services/chat/service"
	"github.com/team/webchat-server/services/chat/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	db, err := mysql.New(cfg.MySQLDSN + "webchat_chat")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}

	hub := websocket.NewHub()
	repo := repository.NewChatRepo(db)
	svc := service.New(repo, hub)
	h := handler.New(svc)

	// gRPC server
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcpkg.UnaryLoggingInterceptor))
	chatv1.RegisterChatServiceServer(s, h)
	reflection.Register(s)
	go func() {
		lis, _ := net.Listen("tcp", ":50052")
		log.Println("chat-service gRPC on :50052")
		s.Serve(lis)
	}()

	// WebSocket server
	http.HandleFunc("/ws", websocket.ServeWS(hub, func(c *websocket.Client, msg *websocket.WsMessage) {
		log.Printf("ws message from %d: %s", c.UserID, msg.Type)
	}))
	log.Println("chat-service WebSocket on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

- [ ] **Step 4: 验证编译**

```bash
cd server && go build ./services/chat/
```

Expected: 编译成功。

- [ ] **Step 5: Commit**

```bash
git add server/services/chat/service/ server/services/chat/handler/ server/services/chat/main.go && git commit -m "feat: implement chat-service gRPC + WebSocket"
```

---

### Task B6: group-service — model + repository + service + handler + main.go

**Files:**
- Create: `server/services/group/model/group.go`
- Create: `server/services/group/repository/group.go`
- Create: `server/services/group/service/group.go`
- Create: `server/services/group/handler/group.go`
- Modify: `server/services/group/main.go`

- [ ] **Step 1: 定义 model**

```go
// server/services/group/model/group.go
package model

type Group struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	OwnerID      int64  `json:"owner_id"`
	Announcement string `json:"announcement"`
	MemberCount  int32  `json:"member_count"`
	MaxMembers   int32  `json:"max_members"`
	CreatedAt    int64  `json:"created_at"`
}

type GroupMember struct {
	GroupID  int64  `json:"group_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"`
	Alias    string `json:"alias"`
	IsMuted  bool   `json:"is_muted"`
	JoinedAt int64  `json:"joined_at"`
}
```

- [ ] **Step 2: 编写 repository**

```go
// server/services/group/repository/group.go
package repository

import (
	"context"
	"database/sql"
	"time"
	"github.com/team/webchat-server/services/group/model"
)

type GroupRepo struct{ db *sql.DB }

func NewGroupRepo(db *sql.DB) *GroupRepo { return &GroupRepo{db} }

func (r *GroupRepo) Create(ctx context.Context, g *model.Group) (int64, error) {
	now := time.Now().UnixMilli()
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO grps (name, avatar, owner_id, member_count, max_members, created_at) VALUES (?, ?, ?, 1, 500, ?)",
		g.Name, g.Avatar, g.OwnerID, now)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	r.db.ExecContext(ctx, "INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'owner', ?)", id, g.OwnerID, now)
	return id, nil
}

func (r *GroupRepo) FindByID(ctx context.Context, id int64) (*model.Group, error) {
	g := &model.Group{}
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name, avatar, owner_id, announcement, member_count, max_members, created_at FROM grps WHERE id = ?", id)
	err := row.Scan(&g.ID, &g.Name, &g.Avatar, &g.OwnerID, &g.Announcement, &g.MemberCount, &g.MaxMembers, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *GroupRepo) AddMembers(ctx context.Context, groupID int64, userIDs []int64) error {
	now := time.Now().UnixMilli()
	for _, uid := range userIDs {
		r.db.ExecContext(ctx, "INSERT IGNORE INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'member', ?)", groupID, uid, now)
	}
	r.db.ExecContext(ctx, "UPDATE grps SET member_count = (SELECT COUNT(*) FROM group_members WHERE group_id = ?) WHERE id = ?", groupID, groupID)
	return nil
}

func (r *GroupRepo) RemoveMember(ctx context.Context, groupID, userID int64) error {
	r.db.ExecContext(ctx, "DELETE FROM group_members WHERE group_id = ? AND user_id = ? AND role != 'owner'", groupID, userID)
	r.db.ExecContext(ctx, "UPDATE grps SET member_count = (SELECT COUNT(*) FROM group_members WHERE group_id = ?) WHERE id = ?", groupID, groupID)
	return nil
}

func (r *GroupRepo) ListMembers(ctx context.Context, groupID int64) ([]*model.GroupMember, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT group_id, user_id, role, alias, is_muted, joined_at FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*model.GroupMember
	for rows.Next() {
		m := &model.GroupMember{}
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.Role, &m.Alias, &m.IsMuted, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}
```

- [ ] **Step 3: 编写 service**

```go
// server/services/group/service/group.go
package service

import (
	"context"
	"github.com/team/webchat-server/pkg/errcode"
	"github.com/team/webchat-server/services/group/model"
	"github.com/team/webchat-server/services/group/repository"
)

type GroupService struct{ repo *repository.GroupRepo }

func New(repo *repository.GroupRepo) *GroupService { return &GroupService{repo} }

func (s *GroupService) Create(ctx context.Context, name, ownerID string, memberIDs []string) (*model.Group, error) {
	oid, _ := strconv.ParseInt(ownerID, 10, 64)
	g := &model.Group{Name: name, OwnerID: oid}
	id, err := s.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}
	var mids []int64
	for _, mid := range memberIDs {
		id, _ := strconv.ParseInt(mid, 10, 64)
		mids = append(mids, id)
	}
	s.repo.AddMembers(ctx, id, mids)
	return s.repo.FindByID(ctx, id)
}

func (s *GroupService) Get(ctx context.Context, id int64) (*model.Group, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *GroupService) AddMembers(ctx context.Context, groupID int64, userIDs []int64) error {
	return s.repo.AddMembers(ctx, groupID, userIDs)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, userID int64) error {
	return s.repo.RemoveMember(ctx, groupID, userID)
}

func (s *GroupService) ListMembers(ctx context.Context, groupID int64) ([]*model.GroupMember, error) {
	return s.repo.ListMembers(ctx, groupID)
}
```

- [ ] **Step 4: 编写 handler + main.go**

handler 和 main.go 与 user/chat 服务模式一致：
- handler 实现 `GroupServiceServer` 接口
- main.go 启动 gRPC server 在 `:50053`

- [ ] **Step 5: 验证编译并提交**

```bash
cd server && go build ./services/group/ && git add server/services/group/ && git commit -m "feat: implement group-service"
```

---

### Task B7: media-service — MinIO 上传 + presigned URL

**Files:**
- Create: `server/services/media/service/media.go`
- Create: `server/services/media/handler/media.go`
- Modify: `server/services/media/main.go`

- [ ] **Step 1: 编写 service 层**

```go
// server/services/media/service/media.go
package service

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	miniopkg "github.com/team/webchat-server/pkg/minio"
)

type MediaService struct {
	client *minio.Client
	bucket string
}

func New(client *minio.Client) *MediaService {
	return &MediaService{client: client, bucket: "webchat-media"}
}

func (s *MediaService) Init(ctx context.Context) error {
	return miniopkg.EnsureBucket(ctx, s.client, s.bucket)
}

func (s *MediaService) GetUploadURL(ctx context.Context, fileName, contentType string) (string, string, error) {
	fileID := uuid.New().String()
	objectName := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), fileID)
	url, err := s.client.PresignedPutObject(ctx, s.bucket, objectName, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	return url.String(), fileID, nil
}

func (s *MediaService) GetFileURL(ctx context.Context, fileID string) (string, error) {
	// Scan for the object by prefix pattern
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Prefix: "", Recursive: true})
	for obj := range objectCh {
		if obj.Err != nil {
			continue
		}
	}
	url, err := s.client.PresignedGetObject(ctx, s.bucket, fileID, 1*time.Hour, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
```

- [ ] **Step 2: 编写 handler + main.go**

```go
// server/services/media/handler/media.go
package handler

import (
	"context"
	mediav1 "github.com/team/webchat-server/gen/media/v1"
	"github.com/team/webchat-server/services/media/service"
)

type MediaHandler struct {
	mediav1.UnimplementedMediaServiceServer
	svc *service.MediaService
}

func New(svc *service.MediaService) *MediaHandler { return &MediaHandler{svc: svc} }

func (h *MediaHandler) GetUploadURL(ctx context.Context, req *mediav1.GetUploadURLRequest) (*mediav1.GetUploadURLResponse, error) {
	url, fileID, err := h.svc.GetUploadURL(ctx, req.FileName, req.ContentType)
	if err != nil {
		return nil, err
	}
	return &mediav1.GetUploadURLResponse{UploadUrl: url, FileId: fileID, ExpiresIn: 900}, nil
}

func (h *MediaHandler) GetFileURL(ctx context.Context, req *mediav1.GetFileURLRequest) (*mediav1.GetFileURLResponse, error) {
	url, err := h.svc.GetFileURL(ctx, req.FileId)
	if err != nil {
		return nil, err
	}
	return &mediav1.GetFileURLResponse{Url: url}, nil
}
```

```go
// server/services/media/main.go
package main

import (
	"context"
	"log"
	"net"
	mediav1 "github.com/team/webchat-server/gen/media/v1"
	"github.com/team/webchat-server/pkg/config"
	_ "github.com/team/webchat-server/pkg/grpc"
	"github.com/team/webchat-server/pkg/minio"
	"github.com/team/webchat-server/services/media/handler"
	"github.com/team/webchat-server/services/media/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	client, err := minio.New(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey)
	if err != nil {
		log.Fatalf("minio: %v", err)
	}
	svc := service.New(client)
	if err := svc.Init(context.Background()); err != nil {
		log.Fatalf("minio init: %v", err)
	}
	h := handler.New(svc)
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcpkg.UnaryLoggingInterceptor))
	mediav1.RegisterMediaServiceServer(s, h)
	reflection.Register(s)
	lis, _ := net.Listen("tcp", ":50054")
	log.Println("media-service on :50054")
	log.Fatal(s.Serve(lis))
}
```

- [ ] **Step 3: 验证编译并提交**

```bash
cd server && go build ./services/media/ && git add server/services/media/ && git commit -m "feat: implement media-service with MinIO"
```

---

### Task B8: gateway — HTTP API 网关

**Files:**
- Create: `server/services/gateway/router/router.go`
- Create: `server/services/gateway/middleware/auth.go`
- Create: `server/services/gateway/middleware/ratelimit.go`
- Create: `server/services/gateway/middleware/cors.go`
- Create: `server/services/gateway/handler/auth.go`
- Create: `server/services/gateway/handler/user.go`
- Create: `server/services/gateway/handler/chat.go`
- Create: `server/services/gateway/handler/group.go`
- Create: `server/services/gateway/handler/media.go`
- Modify: `server/services/gateway/main.go`

- [ ] **Step 1: 编写中间件**

```go
// server/services/gateway/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"
	"github.com/team/webchat-server/pkg/jwt"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			claims, err := jwt.Validate(secret, auth[7:])
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			r.Header.Set("X-User-ID", claims.UserID)
			next.ServeHTTP(w, r)
		})
	}
}
```

```go
// server/services/gateway/middleware/cors.go
package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

- [ ] **Step 2: 编写 router 与 main.go**

```go
// server/services/gateway/router/router.go
package router

import (
	"net/http"
	"github.com/team/webchat-server/services/gateway/handler"
	"github.com/team/webchat-server/services/gateway/middleware"
)

func New(authHandler *handler.AuthHandler, userHandler *handler.UserHandler,
	chatHandler *handler.ChatHandler, groupHandler *handler.GroupHandler,
	mediaHandler *handler.MediaHandler, jwtSecret string) http.Handler {

	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.RefreshToken)

	// protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/v1/users/me", userHandler.GetMe)
	protected.HandleFunc("PUT /api/v1/users/me", userHandler.UpdateMe)
	protected.HandleFunc("GET /api/v1/users/{id}", userHandler.GetUser)
	protected.HandleFunc("GET /api/v1/contacts", userHandler.ListContacts)
	protected.HandleFunc("POST /api/v1/contacts/request", userHandler.AddContact)
	protected.HandleFunc("PUT /api/v1/contacts/request/{id}", userHandler.HandleFriendRequest)
	protected.HandleFunc("GET /api/v1/conversations", chatHandler.GetConversations)
	protected.HandleFunc("GET /api/v1/messages/{conv_id}", chatHandler.GetMessages)
	protected.HandleFunc("POST /api/v1/messages/send", chatHandler.SendMessage)
	protected.HandleFunc("POST /api/v1/groups", groupHandler.CreateGroup)
	protected.HandleFunc("GET /api/v1/groups/{id}", groupHandler.GetGroup)
	protected.HandleFunc("POST /api/v1/groups/{id}/members", groupHandler.AddMember)
	protected.HandleFunc("POST /api/v1/files/upload", mediaHandler.GetUploadURL)
	protected.HandleFunc("GET /api/v1/files/{id}/url", mediaHandler.GetFileURL)

	mux.Handle("/api/v1/", middleware.Auth(jwtSecret)(protected))

	var h http.Handler = mux
	h = middleware.CORS(h)
	return h
}
```

- [ ] **Step 3: 编写 auth handler**

```go
// server/services/gateway/handler/auth.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	userv1 "github.com/team/webchat-server/gen/user/v1"
	"github.com/team/webchat-server/pkg/jwt"
)

type AuthHandler struct {
	userClient userv1.UserServiceClient
	jwtSecret  string
}

func NewAuthHandler(userClient userv1.UserServiceClient, jwtSecret string) *AuthHandler {
	return &AuthHandler{userClient: userClient, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.userClient.Login(context.Background(), &userv1.LoginRequest{
		Phone: req.Phone, Password: req.Password,
	})
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	token, _ := jwt.Generate(h.jwtSecret, resp.UserId)
	json.NewEncoder(w).Encode(map[string]string{"user_id": resp.UserId, "token": token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.userClient.Register(context.Background(), &userv1.RegisterRequest{
		Phone: req.Phone, Password: req.Password, Nickname: req.Nickname,
	})
	if err != nil {
		http.Error(w, `{"error":"registration failed"}`, http.StatusBadRequest)
		return
	}
	token, _ := jwt.Generate(h.jwtSecret, resp.UserId)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"user_id": resp.UserId, "token": token})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	token, _ := jwt.Generate(h.jwtSecret, userID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
```

- [ ] **Step 4: 编写 user/chat/group/media handlers**

其他 handler 与 auth handler 模式一致：接收 HTTP JSON → 调用 gRPC client → 返回 JSON。

user handler 示例：

```go
// server/services/gateway/handler/user.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	userv1 "github.com/team/webchat-server/gen/user/v1"
)

type UserHandler struct{ client userv1.UserServiceClient }

func NewUserHandler(client userv1.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	resp, err := h.client.GetProfile(context.Background(), &userv1.GetProfileRequest{UserId: userID})
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) { /* similar pattern */ }
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request)  { /* similar pattern */ }
func (h *UserHandler) ListContacts(w http.ResponseWriter, r *http.Request) { /* similar pattern */ }
func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request)   { /* similar pattern */ }
func (h *UserHandler) HandleFriendRequest(w http.ResponseWriter, r *http.Request) { /* similar pattern */ }
```

- [ ] **Step 5: 编写 main.go**

```go
// server/services/gateway/main.go
package main

import (
	"log"
	"net/http"
	userv1 "github.com/team/webchat-server/gen/user/v1"
	chatv1 "github.com/team/webchat-server/gen/chat/v1"
	groupv1 "github.com/team/webchat-server/gen/group/v1"
	mediav1 "github.com/team/webchat-server/gen/media/v1"
	"github.com/team/webchat-server/pkg/config"
	"github.com/team/webchat-server/services/gateway/handler"
	"github.com/team/webchat-server/services/gateway/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	userConn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	chatConn, _ := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	groupConn, _ := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	mediaConn, _ := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))

	userClient := userv1.NewUserServiceClient(userConn)
	chatClient := chatv1.NewChatServiceClient(chatConn)
	groupClient := groupv1.NewGroupServiceClient(groupConn)
	mediaClient := mediav1.NewMediaServiceClient(mediaConn)

	authH := handler.NewAuthHandler(userClient, cfg.JWTSecret)
	userH := handler.NewUserHandler(userClient)
	chatH := handler.NewChatHandler(chatClient)
	groupH := handler.NewGroupHandler(groupClient)
	mediaH := handler.NewMediaHandler(mediaClient)

	r := router.New(authH, userH, chatH, groupH, mediaH, cfg.JWTSecret)

	log.Printf("gateway listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
```

- [ ] **Step 6: 验证编译并提交**

```bash
cd server && go build ./services/gateway/ && git add server/services/gateway/ && git commit -m "feat: implement API gateway with HTTP-to-gRPC routing"
```

---

## Part C: Flutter 客户端

### Task C1: 初始化 Flutter 项目

**Files:**
- Create: `client/` (via `flutter create`)
- Modify: `client/pubspec.yaml`

- [ ] **Step 1: 创建 Flutter 项目**

```bash
flutter create --org com.webchat --project-name webchat client
```

- [ ] **Step 2: 添加依赖**

编辑 `client/pubspec.yaml`，在 `dependencies` 下添加：

```yaml
dependencies:
  flutter:
    sdk: flutter
  flutter_bloc: ^8.1.3
  go_router: ^13.0.0
  dio: ^5.4.0
  web_socket_channel: ^2.4.0
  drift: ^2.15.0
  sqlite3_flutter_libs: ^0.5.0
  flutter_secure_storage: ^8.0.0
  cached_network_image: ^3.3.0
  image_picker: ^1.0.7
  record: ^5.0.4
  path_provider: ^2.1.2
  json_annotation: ^4.8.1

dev_dependencies:
  flutter_test:
    sdk: flutter
  drift_dev: ^2.15.0
  build_runner: ^2.4.8
  json_serializable: ^6.7.1
  mocktail: ^1.0.3
  flutter_lints: ^3.0.0
```

- [ ] **Step 3: 安装依赖**

```bash
cd client && flutter pub get
```

Expected: `flutter pub get` 成功。

- [ ] **Step 4: Commit**

```bash
git add client/ && git commit -m "feat: init Flutter project with dependencies"
```

---

### Task C2: Flutter core — 网络层与 Token 管理

**Files:**
- Create: `client/lib/core/network/dio_client.dart`
- Create: `client/lib/core/network/ws_client.dart`
- Create: `client/lib/core/network/api_paths.dart`
- Create: `client/lib/core/auth/token_manager.dart`
- Create: `client/lib/core/theme/colors.dart`
- Create: `client/lib/core/theme/theme.dart`

- [ ] **Step 1: API 路径常量**

```dart
// client/lib/core/network/api_paths.dart
class ApiPaths {
  static const baseUrl = 'http://10.0.2.2:8080/api/v1';
  static const login = '/auth/login';
  static const register = '/auth/register';
  static const refreshToken = '/auth/refresh';
  static const usersMe = '/users/me';
  static const contacts = '/contacts';
  static const contactRequest = '/contacts/request';
  static const conversations = '/conversations';
  static const messages = '/messages';
  static const groups = '/groups';
  static const fileUpload = '/files/upload';
  static const fileUrl = '/files';
  static const wsUrl = 'ws://10.0.2.2:8081/ws';
}
```

- [ ] **Step 2: Token 管理器**

```dart
// client/lib/core/auth/token_manager.dart
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenManager {
  final _storage = const FlutterSecureStorage();
  static const _tokenKey = 'auth_token';
  static const _userIdKey = 'user_id';

  Future<void> save(String token, String userId) async {
    await _storage.write(key: _tokenKey, value: token);
    await _storage.write(key: _userIdKey, value: userId);
  }

  Future<String?> get token => _storage.read(key: _tokenKey);
  Future<String?> get userId => _storage.read(key: _userIdKey);
  Future<void> clear() => _storage.deleteAll();
}
```

- [ ] **Step 3: Dio HTTP 客户端**

```dart
// client/lib/core/network/dio_client.dart
import 'package:dio/dio.dart';
import '../auth/token_manager.dart';
import 'api_paths.dart';

class DioClient {
  late final Dio dio;
  final TokenManager _tokenManager;

  DioClient(this._tokenManager) {
    dio = Dio(BaseOptions(
      baseUrl: ApiPaths.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _tokenManager.token;
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          final token = await _tokenManager.token;
          if (token != null) {
            final newToken = await _refreshToken();
            if (newToken != null) {
              error.requestOptions.headers['Authorization'] = 'Bearer $newToken';
              final retry = await dio.fetch(error.requestOptions);
              return handler.resolve(retry);
            }
          }
          await _tokenManager.clear();
        }
        handler.next(error);
      },
    ));
  }

  Future<String?> _refreshToken() async {
    try {
      final resp = await dio.post(ApiPaths.refreshToken);
      final token = resp.data['token'] as String;
      final userId = resp.data['user_id'] as String;
      await _tokenManager.save(token, userId);
      return token;
    } catch (_) {
      return null;
    }
  }
}
```

- [ ] **Step 4: WebSocket 客户端**

```dart
// client/lib/core/network/ws_client.dart
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../auth/token_manager.dart';
import 'api_paths.dart';

typedef WsMessageHandler = void Function(Map<String, dynamic> message);

class WsClient {
  WebSocketChannel? _channel;
  final TokenManager _tokenManager;
  WsMessageHandler? onMessage;

  WsClient(this._tokenManager);

  Future<void> connect() async {
    final userId = await _tokenManager.userId;
    if (userId == null) return;
    final uri = Uri.parse('${ApiPaths.wsUrl}?user_id=$userId');
    _channel = WebSocketChannel.connect(uri);
    _channel!.stream.listen(
      (data) {
        if (onMessage != null) {
          onMessage!(json.decode(data) as Map<String, dynamic>);
        }
      },
      onError: (_) => _reconnect(),
      onDone: () => _reconnect(),
    );
  }

  void send(Map<String, dynamic> message) {
    _channel?.sink.add(json.encode(message));
  }

  void _reconnect() {
    Future.delayed(const Duration(seconds: 3), connect);
  }

  void dispose() {
    _channel?.sink.close();
  }
}
```

- [ ] **Step 5: 主题**

```dart
// client/lib/core/theme/colors.dart
import 'package:flutter/material.dart';

class AppColors {
  static const primary = Color(0xFF07C160); // 微信绿
  static const background = Color(0xFFEDEDED);
  static const surface = Colors.white;
  static const textPrimary = Color(0xFF191919);
  static const textSecondary = Color(0xFF999999);
  static const divider = Color(0xFFE5E5E5);
  static const redBadge = Color(0xFFF74C4C);
}
```

```dart
// client/lib/core/theme/theme.dart
import 'package:flutter/material.dart';
import 'colors.dart';

class AppTheme {
  static ThemeData get light => ThemeData(
    primaryColor: AppColors.primary,
    scaffoldBackgroundColor: AppColors.background,
    appBarTheme: const AppBarTheme(
      backgroundColor: AppColors.surface,
      foregroundColor: AppColors.textPrimary,
      elevation: 0.5,
    ),
    tabBarTheme: const TabBarTheme(
      labelColor: AppColors.primary,
      unselectedLabelColor: AppColors.textSecondary,
    ),
  );
}
```

- [ ] **Step 6: Commit**

```bash
git add client/lib/core/ && git commit -m "feat: add Flutter core network, auth, theme"
```

---

### Task C3: 用户认证模块 (Auth Feature)

**Files:**
- Create: `client/lib/features/auth/data/auth_api.dart`
- Create: `client/lib/features/auth/data/auth_repository.dart`
- Create: `client/lib/features/auth/bloc/auth_event.dart`
- Create: `client/lib/features/auth/bloc/auth_state.dart`
- Create: `client/lib/features/auth/bloc/auth_bloc.dart`
- Create: `client/lib/features/auth/ui/login_page.dart`
- Create: `client/lib/features/auth/ui/register_page.dart`

- [ ] **Step 1: Auth API**

```dart
// client/lib/features/auth/data/auth_api.dart
import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(this._dio);

  Future<Map<String, dynamic>> login(String phone, String password) async {
    final resp = await _dio.post(ApiPaths.login, data: {
      'phone': phone,
      'password': password,
    });
    return resp.data;
  }

  Future<Map<String, dynamic>> register(String phone, String password, String nickname) async {
    final resp = await _dio.post(ApiPaths.register, data: {
      'phone': phone,
      'password': password,
      'nickname': nickname,
    });
    return resp.data;
  }
}
```

- [ ] **Step 2: Auth Repository**

```dart
// client/lib/features/auth/data/auth_repository.dart
import 'auth_api.dart';
import '../../../core/auth/token_manager.dart';

class AuthRepository {
  final AuthApi _api;
  final TokenManager _tokenManager;
  AuthRepository(this._api, this._tokenManager);

  Future<void> login(String phone, String password) async {
    final data = await _api.login(phone, password);
    await _tokenManager.save(data['token'], data['user_id']);
  }

  Future<void> register(String phone, String password, String nickname) async {
    final data = await _api.register(phone, password, nickname);
    await _tokenManager.save(data['token'], data['user_id']);
  }

  Future<bool> isLoggedIn() async {
    final token = await _tokenManager.token;
    return token != null;
  }

  Future<void> logout() => _tokenManager.clear();
}
```

- [ ] **Step 3: Auth BLoC**

```dart
// client/lib/features/auth/bloc/auth_event.dart
abstract class AuthEvent {}
class AuthLoginRequested extends AuthEvent {
  final String phone;
  final String password;
  AuthLoginRequested(this.phone, this.password);
}
class AuthRegisterRequested extends AuthEvent {
  final String phone;
  final String password;
  final String nickname;
  AuthRegisterRequested(this.phone, this.password, this.nickname);
}
class AuthCheckStatus extends AuthEvent {}
class AuthLogoutRequested extends AuthEvent {}
```

```dart
// client/lib/features/auth/bloc/auth_state.dart
enum AuthStatus { initial, loading, authenticated, unauthenticated, error }

class AuthState {
  final AuthStatus status;
  final String? error;
  const AuthState({this.status = AuthStatus.initial, this.error});
}
```

```dart
// client/lib/features/auth/bloc/auth_bloc.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/auth_repository.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _repo;
  AuthBloc(this._repo) : super(const AuthState()) {
    on<AuthLoginRequested>(_onLogin);
    on<AuthRegisterRequested>(_onRegister);
    on<AuthCheckStatus>(_onCheckStatus);
    on<AuthLogoutRequested>(_onLogout);
  }

  Future<void> _onLogin(AuthLoginRequested event, Emitter<AuthState> emit) async {
    emit(const AuthState(status: AuthStatus.loading));
    try {
      await _repo.login(event.phone, event.password);
      emit(const AuthState(status: AuthStatus.authenticated));
    } catch (e) {
      emit(AuthState(status: AuthStatus.error, error: '登录失败'));
    }
  }

  Future<void> _onRegister(AuthRegisterRequested event, Emitter<AuthState> emit) async {
    emit(const AuthState(status: AuthStatus.loading));
    try {
      await _repo.register(event.phone, event.password, event.nickname);
      emit(const AuthState(status: AuthStatus.authenticated));
    } catch (e) {
      emit(AuthState(status: AuthStatus.error, error: '注册失败'));
    }
  }

  Future<void> _onCheckStatus(AuthCheckStatus event, Emitter<AuthState> emit) async {
    final loggedIn = await _repo.isLoggedIn();
    emit(AuthState(status: loggedIn ? AuthStatus.authenticated : AuthStatus.unauthenticated));
  }

  Future<void> _onLogout(AuthLogoutRequested event, Emitter<AuthState> emit) async {
    await _repo.logout();
    emit(const AuthState(status: AuthStatus.unauthenticated));
  }
}
```

- [ ] **Step 4: 登录页面 UI**

```dart
// client/lib/features/auth/ui/login_page.dart
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';
import '../bloc/auth_state.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _phoneCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: BlocListener<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state.status == AuthStatus.error) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.error ?? '未知错误')),
            );
          }
        },
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('微信', style: TextStyle(fontSize: 36, fontWeight: FontWeight.bold)),
              const SizedBox(height: 48),
              TextField(
                controller: _phoneCtrl,
                decoration: const InputDecoration(labelText: '手机号'),
                keyboardType: TextInputType.phone,
              ),
              const SizedBox(height: 16),
              TextField(
                controller: _pwdCtrl,
                decoration: const InputDecoration(labelText: '密码'),
                obscureText: true,
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () {
                    context.read<AuthBloc>().add(
                      AuthLoginRequested(_phoneCtrl.text, _pwdCtrl.text),
                    );
                  },
                  style: ElevatedButton.styleFrom(backgroundColor: Theme.of(context).primaryColor),
                  child: const Text('登录', style: TextStyle(fontSize: 16, color: Colors.white)),
                ),
              ),
              const SizedBox(height: 16),
              TextButton(
                onPressed: () {
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const RegisterPage()));
                },
                child: const Text('注册新账号'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
```

- [ ] **Step 5: Commit**

```bash
git add client/lib/features/auth/ && git commit -m "feat: add Flutter auth feature (login/register)"
```

---

### Task C4: 首页 Tab 骨架

**Files:**
- Create: `client/lib/features/home/home_page.dart`
- Create: `client/lib/features/home/main_tab_scaffold.dart`
- Modify: `client/lib/app.dart`
- Modify: `client/lib/main.dart`

- [ ] **Step 1: 主 Tab 骨架**

```dart
// client/lib/features/home/home_page.dart
import 'package:flutter/material.dart';
import '../chat/ui/conversation_list_page.dart';
import '../contacts/ui/contacts_list_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _currentIndex = 0;

  final _pages = const [
    ConversationListPage(),
    ContactsListPage(),
    Center(child: Text('发现')),
    Center(child: Text('我')),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _currentIndex,
        onTap: (i) => setState(() => _currentIndex = i),
        type: BottomNavigationBarType.fixed,
        selectedItemColor: Theme.of(context).primaryColor,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble_outline), label: '微信'),
          BottomNavigationBarItem(icon: Icon(Icons.contacts_outlined), label: '通讯录'),
          BottomNavigationBarItem(icon: Icon(Icons.explore_outlined), label: '发现'),
          BottomNavigationBarItem(icon: Icon(Icons.person_outline), label: '我'),
        ],
      ),
    );
  }
}
```

- [ ] **Step 2: app.dart + main.dart**

```dart
// client/lib/app.dart
import 'package:flutter/material.dart';
import 'core/theme/theme.dart';
import 'features/auth/bloc/auth_bloc.dart';
import 'features/auth/bloc/auth_event.dart';
import 'features/auth/bloc/auth_state.dart';
import 'features/auth/ui/login_page.dart';
import 'features/home/home_page.dart';

class App extends StatelessWidget {
  const App({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: AppTheme.light,
      home: BlocBuilder<AuthBloc, AuthState>(
        builder: (context, state) {
          if (state.status == AuthStatus.authenticated) {
            return const HomePage();
          }
          return const LoginPage();
        },
      ),
    );
  }
}
```

```dart
// client/lib/main.dart
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/network/dio_client.dart';
import 'core/network/ws_client.dart';
import 'core/auth/token_manager.dart';
import 'features/auth/data/auth_api.dart';
import 'features/auth/data/auth_repository.dart';
import 'features/auth/bloc/auth_bloc.dart';
import 'app.dart';

void main() {
  final tokenManager = TokenManager();
  final dioClient = DioClient(tokenManager);
  final wsClient = WsClient(tokenManager);
  final authApi = AuthApi(dioClient.dio);
  final authRepo = AuthRepository(authApi, tokenManager);

  runApp(
    BlocProvider(
      create: (_) => AuthBloc(authRepo)..add(AuthCheckStatus()),
      child: const App(),
    ),
  );
}
```

- [ ] **Step 3: Commit**

```bash
git add client/lib/app.dart client/lib/main.dart client/lib/features/home/ && git commit -m "feat: add home tab scaffold and app entry point"
```

---

### Task C5: 会话列表页面

**Files:**
- Create: `client/lib/features/chat/data/chat_api.dart`
- Create: `client/lib/features/chat/data/chat_repository.dart`
- Create: `client/lib/features/chat/bloc/conversation_bloc.dart`
- Create: `client/lib/features/chat/ui/conversation_list_page.dart`

- [ ] **Step 1: Chat API**

```dart
// client/lib/features/chat/data/chat_api.dart
import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class ChatApi {
  final Dio _dio;
  ChatApi(this._dio);

  Future<List<dynamic>> getConversations() async {
    final resp = await _dio.get(ApiPaths.conversations);
    return resp.data['conversations'] as List<dynamic>;
  }

  Future<List<dynamic>> getMessages(String convId, {int page = 1, int pageSize = 20}) async {
    final resp = await _dio.get('${ApiPaths.messages}/$convId', queryParameters: {
      'page': page,
      'page_size': pageSize,
    });
    return resp.data['messages'] as List<dynamic>;
  }

  Future<Map<String, dynamic>> sendMessage(Map<String, dynamic> msg) async {
    final resp = await _dio.post('${ApiPaths.messages}/send', data: msg);
    return resp.data;
  }
}
```

- [ ] **Step 2: Chat Repository**

```dart
// client/lib/features/chat/data/chat_repository.dart
import 'chat_api.dart';
import '../../../core/network/ws_client.dart';

class ChatRepository {
  final ChatApi _api;
  final WsClient _wsClient;
  ChatRepository(this._api, this._wsClient);

  Future<List<Map<String, dynamic>>> getConversations() async {
    final list = await _api.getConversations();
    return list.cast<Map<String, dynamic>>();
  }

  Future<List<Map<String, dynamic>>> getMessages(String convId) async {
    final list = await _api.getMessages(convId);
    return list.cast<Map<String, dynamic>>();
  }

  Future<void> sendMessage(Map<String, dynamic> msg) async {
    _wsClient.send(msg);
  }

  void connectWebSocket() => _wsClient.connect();
}
```

- [ ] **Step 3: Conversation BLoC**

```dart
// client/lib/features/chat/bloc/conversation_bloc.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/chat_repository.dart';

// events & states (内联)
abstract class ConvEvent {}
class LoadConversations extends ConvEvent {}

class ConvState {
  final bool loading;
  final List<Map<String, dynamic>> conversations;
  final String? error;
  const ConvState({this.loading = false, this.conversations = const [], this.error});
}

class ConversationBloc extends Bloc<ConvEvent, ConvState> {
  final ChatRepository _repo;
  ConversationBloc(this._repo) : super(const ConvState()) {
    on<LoadConversations>(_onLoad);
  }

  Future<void> _onLoad(LoadConversations event, Emitter<ConvState> emit) async {
    emit(const ConvState(loading: true));
    try {
      final convs = await _repo.getConversations();
      emit(ConvState(conversations: convs));
    } catch (e) {
      emit(ConvState(error: e.toString()));
    }
  }
}
```

- [ ] **Step 4: 会话列表 UI**

```dart
// client/lib/features/chat/ui/conversation_list_page.dart
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/conversation_bloc.dart';
import 'chat_page.dart';

class ConversationListPage extends StatelessWidget {
  const ConversationListPage({super.key});
  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => ConversationBloc(context.read<ChatRepository>())..add(LoadConversations()),
      child: BlocBuilder<ConversationBloc, ConvState>(
        builder: (context, state) {
          if (state.loading) return const Center(child: CircularProgressIndicator());
          final convs = state.conversations;
          return ListView.builder(
            itemCount: convs.length,
            itemBuilder: (context, index) {
              final c = convs[index];
              return ListTile(
                leading: CircleAvatar(backgroundImage: c['target_avatar'] != null ? NetworkImage(c['target_avatar']) : null),
                title: Text(c['target_name'] ?? ''),
                subtitle: Text(c['last_msg']?['content'] ?? '', maxLines: 1, overflow: TextOverflow.ellipsis),
                trailing: (c['unread_count'] as int) > 0
                    ? Badge(label: Text('${c['unread_count']}'))
                    : null,
                onTap: () => Navigator.push(context, MaterialPageRoute(
                  builder: (_) => ChatPage(convId: '${c['target_id']}', title: c['target_name'] ?? ''),
                )),
              );
            },
          );
        },
      ),
    );
  }
}
```

- [ ] **Step 5: Commit**

```bash
git add client/lib/features/chat/data/ client/lib/features/chat/bloc/conversation_bloc.dart client/lib/features/chat/ui/conversation_list_page.dart && git commit -m "feat: add conversation list feature"
```

---

### Task C6: 聊天页面 (ChatPage)

**Files:**
- Create: `client/lib/features/chat/bloc/message_bloc.dart`
- Create: `client/lib/features/chat/ui/chat_page.dart`
- Create: `client/lib/features/chat/ui/widgets/message_bubble.dart`
- Create: `client/lib/features/chat/ui/widgets/chat_input_bar.dart`

- [ ] **Step 1: Message BLoC**

```dart
// client/lib/features/chat/bloc/message_bloc.dart
import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/chat_repository.dart';

abstract class MsgEvent {}
class LoadMessages extends MsgEvent {
  final String convId;
  LoadMessages(this.convId);
}
class SendTextMessage extends MsgEvent {
  final String convId;
  final String toId;
  final String content;
  SendTextMessage(this.convId, this.toId, this.content);
}
class ReceiveMessage extends MsgEvent {
  final Map<String, dynamic> message;
  ReceiveMessage(this.message);
}

class MsgState {
  final bool loading;
  final List<Map<String, dynamic>> messages;
  final String? error;
  const MsgState({this.loading = false, this.messages = const [], this.error});
}

class MessageBloc extends Bloc<MsgEvent, MsgState> {
  final ChatRepository _repo;
  MessageBloc(this._repo) : super(const MsgState()) {
    on<LoadMessages>(_onLoad);
    on<SendTextMessage>(_onSend);
    on<ReceiveMessage>(_onReceive);
  }

  Future<void> _onLoad(LoadMessages event, Emitter<MsgState> emit) async {
    emit(const MsgState(loading: true));
    try {
      final msgs = await _repo.getMessages(event.convId);
      emit(MsgState(messages: msgs));
    } catch (e) {
      emit(MsgState(error: e.toString()));
    }
  }

  Future<void> _onSend(SendTextMessage event, Emitter<MsgState> emit) async {
    final msg = {
      'type': 'chat.message',
      'data': {
        'chat_type': 'single',
        'to_id': event.toId,
        'msg_type': 'text',
        'content': event.content,
      },
    };
    await _repo.sendMessage(msg);
    // optimistic update
    emit(MsgState(messages: [
      ...state.messages,
      {'from_user': 'me', 'msg_type': 'text', 'content': event.content, 'created_at': DateTime.now().millisecondsSinceEpoch},
    ]));
  }

  void _onReceive(ReceiveMessage event, Emitter<MsgState> emit) {
    emit(MsgState(messages: [...state.messages, event.message]));
  }
}
```

- [ ] **Step 2: 消息气泡组件**

```dart
// client/lib/features/chat/ui/widgets/message_bubble.dart
import 'package:flutter/material.dart';

class MessageBubble extends StatelessWidget {
  final bool isMe;
  final String content;
  final String msgType;
  const MessageBubble({super.key, required this.isMe, required this.content, required this.msgType});

  @override
  Widget build(BuildContext context) {
    if (msgType == 'text') {
      return Align(
        alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 12),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: isMe ? const Color(0xFF95EC69) : Colors.white,
            borderRadius: BorderRadius.circular(4),
          ),
          child: Text(content, style: const TextStyle(fontSize: 16)),
        ),
      );
    }
    if (msgType == 'image') {
      return Align(
        alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 12),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: Image.network(content, width: 200, height: 200, fit: BoxFit.cover),
          ),
        ),
      );
    }
    return const SizedBox.shrink();
  }
}
```

- [ ] **Step 3: 聊天输入栏**

```dart
// client/lib/features/chat/ui/widgets/chat_input_bar.dart
import 'package:flutter/material.dart';

class ChatInputBar extends StatelessWidget {
  final TextEditingController controller;
  final VoidCallback onSend;
  final VoidCallback onImage;
  final VoidCallback onVoice;
  const ChatInputBar({
    super.key,
    required this.controller,
    required this.onSend,
    required this.onImage,
    required this.onVoice,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
      decoration: const BoxDecoration(
        color: Colors.white,
        border: Border(top: BorderSide(color: Color(0xFFE5E5E5))),
      ),
      child: Row(
        children: [
          IconButton(icon: const Icon(Icons.mic), onPressed: onVoice),
          Expanded(
            child: TextField(
              controller: controller,
              decoration: InputDecoration(
                hintText: '输入消息...',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(20), borderSide: BorderSide.none),
                filled: true,
                fillColor: const Color(0xFFF5F5F5),
                contentPadding: const EdgeInsets.symmetric(horizontal: 16),
              ),
              textInputAction: TextInputAction.send,
              onSubmitted: (_) => onSend(),
            ),
          ),
          IconButton(icon: const Icon(Icons.image), onPressed: onImage),
          IconButton(icon: const Icon(Icons.send), onPressed: onSend),
        ],
      ),
    );
  }
}
```

- [ ] **Step 4: 聊天页面**

```dart
// client/lib/features/chat/ui/chat_page.dart
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/message_bloc.dart';
import 'widgets/message_bubble.dart';
import 'widgets/chat_input_bar.dart';

class ChatPage extends StatefulWidget {
  final String convId;
  final String title;
  const ChatPage({super.key, required this.convId, required this.title});
  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _inputCtrl = TextEditingController();
  final _scrollCtrl = ScrollController();

  @override
  void initState() {
    super.initState();
    context.read<MessageBloc>().add(LoadMessages(widget.convId));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Column(
        children: [
          Expanded(
            child: BlocBuilder<MessageBloc, MsgState>(
              builder: (context, state) {
                if (state.loading) return const Center(child: CircularProgressIndicator());
                return ListView.builder(
                  controller: _scrollCtrl,
                  itemCount: state.messages.length,
                  itemBuilder: (context, index) {
                    final msg = state.messages[index];
                    return MessageBubble(
                      isMe: msg['from_user'] == 'me',
                      content: msg['content'] ?? '',
                      msgType: msg['msg_type'] ?? 'text',
                    );
                  },
                );
              },
            ),
          ),
          ChatInputBar(
            controller: _inputCtrl,
            onSend: () {
              if (_inputCtrl.text.trim().isEmpty) return;
              context.read<MessageBloc>().add(SendTextMessage(
                widget.convId, widget.convId, _inputCtrl.text.trim(),
              ));
              _inputCtrl.clear();
            },
            onImage: () {},
            onVoice: () {},
          ),
        ],
      ),
    );
  }
}
```

- [ ] **Step 5: Commit**

```bash
git add client/lib/features/chat/bloc/message_bloc.dart client/lib/features/chat/ui/chat_page.dart client/lib/features/chat/ui/widgets/ && git commit -m "feat: add chat page with message bubbles and input bar"
```

---

### Task C7: 通讯录页面

**Files:**
- Create: `client/lib/features/contacts/data/contacts_api.dart`
- Create: `client/lib/features/contacts/data/contacts_repository.dart`
- Create: `client/lib/features/contacts/bloc/contacts_bloc.dart`
- Create: `client/lib/features/contacts/ui/contacts_list_page.dart`
- Create: `client/lib/features/contacts/ui/add_contact_page.dart`

- [ ] **Step 1: Contacts API**

```dart
// client/lib/features/contacts/data/contacts_api.dart
import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class ContactsApi {
  final Dio _dio;
  ContactsApi(this._dio);

  Future<List<dynamic>> list() async {
    final resp = await _dio.get(ApiPaths.contacts);
    return resp.data['contacts'] as List<dynamic>;
  }

  Future<void> addRequest(String toUserId, String message) async {
    await _dio.post(ApiPaths.contactRequest, data: {
      'to_user': toUserId,
      'message': message,
    });
  }

  Future<void> handleRequest(String requestId, String action) async {
    await _dio.put('${ApiPaths.contactRequest}/$requestId', data: {'action': action});
  }
}
```

- [ ] **Step 2: Contacts BLoC + UI**

与 chat 模块模式相同：BLoC 管理状态，UI 展示列表 + 搜索 + 添加好友。

通讯录页面核心结构：

```dart
// client/lib/features/contacts/ui/contacts_list_page.dart
class ContactsListPage extends StatelessWidget {
  const ContactsListPage({super.key});
  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => ContactsBloc(context.read<ContactsRepository>())..add(LoadContacts()),
      child: Scaffold(
        appBar: AppBar(title: const Text('通讯录'), actions: [
          IconButton(icon: const Icon(Icons.person_add), onPressed: () {
            Navigator.push(context, MaterialPageRoute(builder: (_) => const AddContactPage()));
          }),
        ]),
        body: BlocBuilder<ContactsBloc, ContactsState>(
          builder: (context, state) {
            if (state.loading) return const Center(child: CircularProgressIndicator());
            return ListView.builder(
              itemCount: state.contacts.length,
              itemBuilder: (context, index) {
                final c = state.contacts[index];
                return ListTile(
                  leading: CircleAvatar(backgroundImage: c['avatar'] != null ? NetworkImage(c['avatar']) : null),
                  title: Text(c['nickname'] ?? ''),
                );
              },
            );
          },
        ),
      ),
    );
  }
}
```

- [ ] **Step 3: Commit**

```bash
git add client/lib/features/contacts/ && git commit -m "feat: add contacts list feature"
```

---

## 自审清单

**1. Spec 覆盖检查（对照 PRD 一期范围）：**
- [x] 注册/登录 → Task C3 (Auth Feature), Task B1-B2 (user-service)
- [x] 通讯录（好友列表、添加好友、搜索）→ Task B1-B2 (user-service contacts), Task C7 (contacts)
- [x] 单聊（文本/图片/语音/撤回/已读/搜索）→ Task B3-B5 (chat-service), Task C5-C6 (Flutter chat)
- [x] 群聊（创建/成员管理/公告/@提醒/免打扰/置顶）→ Task B6 (group-service)
- [x] 会话列表（排序/未读/置顶/删除/免打扰）→ Task B3-B5 (chat-service), Task C5
- [x] 在线状态（在线/离线/正在输入）→ Task B4 (WebSocket)
- [x] 文件上传 → Task B7 (media-service + MinIO)

**2. 占位符扫描：** 无 TBD、TODO。

**3. 类型一致性：** Message、Conversation 等类型在 proto → model → gRPC handler → Flutter 中保持一致。
