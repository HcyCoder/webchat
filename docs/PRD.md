# 仿微信 Android APP — 产品需求文档 (PRD)

## 1. 项目概述

### 1.1 项目背景

开发一款仿微信的安卓即时通讯应用，完整覆盖微信核心功能，包括即时通讯、社交动态、音视频通话、支付和小程序等。项目采用全栈自研模式，后端与客户端均由团队自主开发。

### 1.2 项目定位

面向 C 端用户的超级社交应用，提供一站式沟通与生活服务。

### 1.3 目标平台

- **主平台**：Android（Flutter 跨平台，保留 iOS 扩展能力）

---

## 2. 功能范围

### 2.1 一期：IM 核心

| 模块 | 功能点 |
|------|--------|
| 注册/登录 | 手机号注册、密码登录、短信验证码登录、Token 自动刷新 |
| 通讯录 | 好友列表、好友搜索、添加好友（申请/同意/拒绝）、备注名、标签、拉黑 |
| 单聊 | 文本消息、图片消息（拍照/相册）、语音消息（录制/播放）、表情、消息撤回、消息已读/未读回执、聊天记录搜索 |
| 群聊 | 创建群聊、群成员管理、群公告、@提醒、群聊免打扰/置顶 |
| 会话列表 | 会话列表（按时间排序）、未读角标、置顶、删除会话、免打扰 |
| 在线状态 | 好友在线/离线状态、正在输入提示 |

### 2.2 二期：社交系统

| 模块 | 功能点 |
|------|--------|
| 朋友圈 | 发布动态（文字/图片/视频）、好友时间线、点赞、评论（含回复评论）、可见权限控制 |
| 个人主页 | 头像/昵称/微信号/地区/个性签名、我的朋友圈相册 |
| 发现页 | 朋友圈入口、扫一扫入口 |

### 2.3 三期：扩展系统

| 模块 | 功能点 |
|------|--------|
| 音视频通话 | 一对一语音通话、一对一视频通话、通话记录、WebRTC 实现 |
| 红包 | 普通红包、拼手气红包、拆红包、红包记录、红包过期退款 |
| 转账 | 好友转账、转账记录、钱包余额 |
| 小程序 | 小程序加载容器（WebView + JS Bridge）、小程序管理 |
| 推送 | 离线消息推送、好友请求推送、朋友圈互动推送 |
| 设置 | 账号安全、隐私设置、通用设置、存储管理 |

---

## 3. 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 客户端框架 | Flutter 3.x | Dart 语言，跨平台 |
| 状态管理 | flutter_bloc | BLoC 模式 |
| 路由 | go_router | 深层链接支持 |
| HTTP 客户端 | dio | 拦截器、Token 刷新 |
| 本地存储 | drift + flutter_secure_storage | 结构化数据 + 敏感数据 |
| 后端框架 | Go + go-zero | 微服务框架（API Gateway + RPC） |
| API 协议 | go-zero API (HTTP) + gRPC (Protobuf) | API Gateway 用 .api 定义，服务间用 proto |
| 代码生成 | goctl | 从 .api / .proto 自动生成 handler、client、server |
| 服务间通信 | gRPC (Protobuf) | 同步调用，go-zero RPC client 内置服务发现 |
| 服务发现 | etcd | go-zero 内置 etcd 服务注册与发现 |
| 消息队列 | Apache Kafka | 异步事件总线 |
| 认证方式 | Token + Redis | SHA256 生成 token，Redis 存储并设置 TTL |
| 关系型数据库 | MySQL 8.x | 每服务独立数据库 |
| 缓存 | Redis Cluster | Token、在线状态、消息缓存 |
| 对象存储 | MinIO | 图片、语音、视频、文件 |
| 容器编排 | Kubernetes | 服务部署与伸缩 |
| CI/CD | GitLab CI / GitHub Actions | 自动化构建部署 |
| 日志 | ELK (Filebeat → ES → Kibana) | 集中日志 |
| 监控 | Prometheus + Grafana | 指标采集与告警 |
| 链路追踪 | Jaeger + OpenTelemetry | 分布式追踪 |
| 配置管理 | K8s ConfigMap + Secrets | 环境配置 |

---

## 4. 系统架构

### 4.1 微服务拆分

| 服务 | 类型 | 职责 | 通信方式 |
|------|------|------|----------|
| gateway | go-zero API | API 网关，统一鉴权（查 Redis token→user_id）、限流、路由转发到下游 RPC | HTTP + 调用 RPC client |
| user-service | go-zero RPC | 注册/登录/Token 管理、用户信息、通讯录/好友管理 | gRPC |
| chat-service | go-zero RPC | 单聊/群聊消息收发、消息存储、已读回执 | gRPC + WebSocket（独立 HTTP 端口） |
| group-service | go-zero RPC | 群创建/管理、群成员管理、群公告 | gRPC |
| social-service | go-zero RPC | 朋友圈动态发布/浏览、点赞评论、个人主页 | gRPC |
| media-service | go-zero RPC | 图片/语音/视频/文件上传、转码、CDN 分发 | gRPC + HTTP（上传预签名 URL） |
| call-service | go-zero RPC | 音视频通话信令、房间管理 (WebRTC) | gRPC + WebSocket（独立 HTTP 端口） |
| payment-service | go-zero RPC | 红包收发、转账、钱包 | gRPC |
| miniapp-service | go-zero RPC | 小程序管理、加载、运行沙箱 | gRPC |
| notification-service | go-zero RPC | 推送通知、消息提醒 | gRPC |

所有 RPC 服务启动时自动注册到 etcd，gateway 通过 go-zero RPC client 自动发现并调用下游服务。

### 4.2 部署架构

```
                      ┌──────────────┐
                      │  Nginx LB    │
                      │  (HTTPS)     │
                      └──────┬───────┘
                             │
                      ┌──────┴───────┐
                      │  Gateway     │  (go-zero API, K8s 2+ 副本)
                      │  (Token 鉴权) │  ← 查 Redis 获取 user_id
                      └──────┬───────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
    ┌─────┴─────┐      ┌─────┴─────┐     ┌─────┴─────┐
    │   etcd    │      │ WebSocket │     │   Kafka   │
    │ 服务注册   │      │ 直连      │     │  Cluster  │
    │ 服务发现   │      └─────┬─────┘     └─────┬─────┘
    └─────┬─────┘            │                  │
          │                  │                  │
    ┌─────┴──────────────────┴──────────────────┴──────────┐
    │                                                       │
    │  ┌──────────┐ ┌──────────┐ ┌──────────┐             │
    │  │chat-svc  │ │social-svc│ │group-svc │  ...        │
    │  │MySQL     │ │MySQL     │ │MySQL     │             │
    │  │WebSocket │ └──────────┘ └──────────┘             │
    │  └──────────┘                                        │
    │                                                       │
    │  ┌──────────┐                                        │
    │  │media-svc │── MinIO                                │
    │  └──────────┘                                        │
    │                                                       │
    │  ┌─────────────────────────────────────┐             │
    │  │          Redis Cluster              │             │
    │  │  Token / 在线状态 / 消息缓存         │             │
    │  └─────────────────────────────────────┘             │
    │                                                       │
    │  ┌─────────────────────────────────────┐             │
    │  │          Kafka Cluster              │             │
    │  └─────────────────────────────────────┘             │
    └───────────────────────────────────────────────────────┘
```

每个服务独立 HPA 弹性伸缩，chat-service 和 call-service 使用 StatefulSet 保证 WebSocket 会话粘滞。所有 RPC 服务启动时自动向 etcd 注册，gateway 通过 etcd 自动发现下游服务。

### 4.3 Flutter 客户端架构

```
lib/
├── core/                    # 公共基础设施
│   ├── network/            # WebSocket 长连接、HTTP 请求封装
│   ├── cache/              # 本地缓存
│   ├── storage/            # 本地文件存储
│   ├── auth/               # Token 管理、自动刷新
│   └── theme/              # 主题、国际化
├── features/               # 业务模块（按微信 Tab 划分）
│   ├── chat/               # 消息列表 + 单聊/群聊
│   ├── contacts/           # 通讯录
│   ├── discover/           # 发现页
│   ├── moments/            # 朋友圈
│   ├── me/                 # 我（个人主页、钱包、设置）
│   ├── call/               # 音视频通话
│   └── miniapp/            # 小程序容器
└── shared/                 # 跨 feature 共享组件
```

| 能力 | 选型 |
|------|------|
| 状态管理 | flutter_bloc |
| 路由 | go_router |
| 网络 | dio + WebSocket 长连接 |
| 本地存储 | drift (结构化) + flutter_secure_storage (Token) |
| 图片加载 | cached_network_image |
| WebRTC | flutter_webrtc |
| 音视频录制 | record + camera |
| 推送 | firebase_messaging / 厂商 SDK 桥接 |

---

## 5. 核心数据模型

### 5.1 user-service

- **users** — id, phone, password_hash, nickname, avatar, gender, region, signature, created_at
- **contacts** — user_id, contact_id, remark, tag, is_blocked, added_at
- **friend_requests** — from_user, to_user, message, status(pending/accepted/rejected), created_at

### 5.2 chat-service

- **messages** — id, chat_type(single/group), from_user, to_id, msg_type(text/image/voice/video/file/location), content, is_recalled, created_at
- **conversations** — user_id, chat_type, target_id, last_msg_id, unread_count, is_pinned, is_muted, updated_at
- **read_receipts** — msg_id, user_id, read_at

### 5.3 group-service

- **groups** — id, name, avatar, owner_id, announcement, member_count, max_members, created_at
- **group_members** — group_id, user_id, role(owner/admin/member), alias, is_muted, joined_at

### 5.4 social-service

- **moments** — id, user_id, content_type(text/image/video), content, location, permission, created_at
- **moment_likes** — moment_id, user_id, created_at
- **moment_comments** — moment_id, user_id, reply_to_id, content, created_at

### 5.5 payment-service

- **wallets** — user_id, balance(fen), pay_password_hash, created_at
- **transactions** — id, from_user, to_user, amount, type(red_packet/transfer), status, created_at
- **red_packets** — id, sender_id, total_amount, count, blessing, type(random/normal), expire_at

---

## 6. API 设计

### 6.1 HTTP REST API（Flutter ↔ Gateway）

```
# 认证
POST   /api/v1/auth/login        # 登录，返回 token + user_id
POST   /api/v1/auth/sms-code     # 短信验证码登录
POST   /api/v1/auth/register     # 注册，返回 token + user_id
POST   /api/v1/auth/refresh      # 刷新 token（旧 token 保留 5 分钟后失效）

# 用户
GET    /api/v1/users/me
PUT    /api/v1/users/me
GET    /api/v1/users/{id}

# 通讯录
GET    /api/v1/contacts
POST   /api/v1/contacts/request
PUT    /api/v1/contacts/request/{id}

# 会话与消息
GET    /api/v1/conversations
GET    /api/v1/messages/{conv_id}
POST   /api/v1/messages/send

# 群聊
POST   /api/v1/groups
GET    /api/v1/groups/{id}
POST   /api/v1/groups/{id}/members

# 朋友圈
GET    /api/v1/moments
POST   /api/v1/moments
POST   /api/v1/moments/{id}/like
POST   /api/v1/moments/{id}/comment

# 文件
POST   /api/v1/files/upload
GET    /api/v1/files/{id}/url

# 支付
POST   /api/v1/redpackets/send
POST   /api/v1/redpackets/{id}/open
POST   /api/v1/transfers
```

### 6.2 WebSocket 协议（Flutter ↔ chat-service / call-service）

消息格式：`{ "type": "...", "data": {...}, "seq": 123 }`

| type | 方向 | 说明 |
|------|------|------|
| `chat.message` | 双向 | 发送/接收消息 |
| `chat.recall` | 双向 | 撤回消息 |
| `chat.read` | 双向 | 已读回执 |
| `chat.typing` | 双向 | 对方正在输入 |
| `presence.online` | 服务端推送 | 好友上线 |
| `presence.offline` | 服务端推送 | 好友下线 |
| `call.invite` | 双向 | 通话邀请 |
| `call.answer` | 双向 | 接听 |
| `call.reject` | 双向 | 拒接 |
| `call.sdp` | 双向 | WebRTC SDP 交换 |
| `call.ice` | 双向 | WebRTC ICE Candidate |
| `call.hangup` | 双向 | 挂断 |

### 6.3 内部服务通信

- 同步调用：gRPC (Protobuf)
- 异步事件：Kafka

---

## 7. 非功能性需求

### 7.1 性能指标

| 指标 | 目标 |
|------|------|
| 消息发送延迟 | < 100ms (在线用户) |
| 图片加载首帧 | < 500ms |
| 朋友圈首屏加载 | < 1s |
| 音视频通话延迟 | < 200ms (端到端) |
| 接口响应时间 P99 | < 500ms |

### 7.2 可用性

- 核心服务（chat/gateway/user）可用性 ≥ 99.9%
- 非核心服务可用性 ≥ 99%
- 消息不丢、不重（at-least-once + 客户端去重）
- 数据多副本 + 定期备份

### 7.3 安全

- 全链路 HTTPS + WSS
- 密码 bcrypt 哈希，支付密码独立哈希
- Token + Redis 认证：`SHA256(user_id + timestamp + random)` 生成 token，Redis 存储 `token:{token}` → user_id，TTL 24h
- 敏感数据（消息内容）可落盘加密
- 接口防刷（go-zero 内置限流 + 验证码）
- WebSocket 连接鉴权（通过 token 参数）

### 7.4 可扩展性

- 消息表可按时间/用户维度分片
- chat-service 可按用户 ID 哈希无状态水平扩展
- 朋友圈时间线可用 Redis 缓存热点数据
- 小程序容器独立部署，不耦合核心服务

---

## 8. 开发分期

| 阶段 | 内容 | 预计周期 |
|------|------|----------|
| 一期 | 微服务基础设施 (gateway + 服务注册 + CI/CD)，user-service + chat-service + group-service + media-service + Flutter 聊天核心 | 2-3 月 |
| 二期 | social-service + Flutter 朋友圈 & 个人主页 | 1-2 月 |
| 三期 | call-service + payment-service + miniapp-service + notification-service + Flutter 剩余功能 | 2-3 月 |
| 四期 | 压测优化、安全加固、上线 | 1 月 |

---

## 9. 风险与应对

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| WebSocket 长连接大规模并发 | chat-service 过载 | 连接横向扩展 + 会话粘滞，做好连接数压测 |
| 消息丢失 | 用户体验差 | Kafka 持久化 + ACK 机制 + 客户端去重 |
| 音视频通话质量差 | 功能不可用 | WebRTC TURN/STUN 服务器部署，提前测试弱网 |
| 微服务过多导致运维复杂 | 交付延期 | 初期监控/日志/CI/CD 一次性搭好，降低后续运维成本 |
| MinIO 单点故障 | 文件不可用 | MinIO 集群模式 + 纠删码 |
| 支付安全风险 | 资损 | 支付独立安全审计，金额用分存储避免浮点，操作日志完整记录 |
