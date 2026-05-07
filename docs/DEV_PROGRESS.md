# 仿微信 Android APP — 开发进度

> 最后更新：2026-05-07

## 项目文件清单

| 文件 | 说明 |
|------|------|
| `docs/PRD.md` | 产品需求文档 |
| `docs/superpowers/plans/2026-05-07-phase1-im-core.md` | 一期实施计划 |
| `docs/DEV_PROGRESS.md` | 本文档，开发进度追踪 |
| `server/docker-compose.yml` | 本地开发中间件 |

---

## 一、环境搭建 ✅

| 组件 | 版本 | 状态 | 备注 |
|------|------|------|------|
| Go | 1.26 | ✅ | |
| Docker | 27.3 | ✅ | |
| Docker Compose | v2.29 | ✅ | |
| Flutter | 3.35.7 | ✅ | `/home/hcy/flutter/bin` |
| JDK | 17 (Temurin) | ✅ | `/home/hcy/.local/jdk` |
| Android SDK | 36.1 | ✅ | `/home/hcy/android-sdk` |

### 中间件 (docker compose)

```bash
cd /home/hcy/workspace/webchat/server && docker compose up -d
```

| 服务 | 端口 | 账号/密码 |
|------|------|-----------|
| MySQL 8.0 | 3306 | root / root123 |
| Redis 7 | 6379 | — |
| MinIO | 9000 (API) / 9001 (Console) | minioadmin / minioadmin |
| Kafka | 9092 | — |
| etcd | 2379 | — |

### 环境变量

```bash
export PATH="$PATH:/home/hcy/flutter/bin"
export JAVA_HOME=/home/hcy/.local/jdk
export ANDROID_HOME=/home/hcy/android-sdk
export PATH="$JAVA_HOME/bin:$ANDROID_HOME/platform-tools:$PATH"
```

> 已写入 `~/.bashrc`，新终端自动生效。当前终端需手动 `source ~/.bashrc`。

---

## 二、一期实施进度

一期分为 3 部分共 20 个任务，详见实施计划文档。

### Part A：开发环境与基础设施

| 任务 | 内容 | 状态 |
|------|------|------|
| A1 | 初始化 go-zero 项目骨架 | ✅ 已完成 |
| A2 | 编写 docker-compose.yml | ✅ 已完成 |
| A3 | 定义 Protobuf + API 文件 + goctl 生成 | ✅ 已完成 |
| A4 | 创建公共包 (token/mysql/errcode) | ✅ 已完成 |
| A5 | 数据库迁移文件 | ✅ 已完成 |

### Part B：后端服务

| 任务 | 内容 | 状态 |
|------|------|------|
| B1 | user-service — 配置 + ServiceContext + DAO | ✅ 已完成 |
| B2 | user-service — Logic 实现 | ✅ 已完成 |
| B3 | chat-service — 配置 + DAO + WebSocket | ✅ 已完成 |
| B4 | chat-service — Logic 实现 | ✅ 已完成 |
| B5 | group-service — 完整实现 | ✅ 已完成 |
| B6 | media-service — 完整实现 | ✅ 已完成 |
| B7 | gateway — 中间件 + Logic 对接 RPC | ✅ 已完成 |

### Part C：Flutter 客户端

| 任务 | 内容 | 状态 |
|------|------|------|
| C1 | 初始化 Flutter 项目 | ✅ 已完成 |
| C2 | core — 网络层 + Token 管理 + 主题 | ⬜ 未开始 |
| C3 | Auth Feature — 登录/注册 | ⬜ 未开始 |
| C4 | 首页 Tab 骨架 | ⬜ 未开始 |
| C5 | 会话列表页面 | ⬜ 未开始 |
| C6 | 聊天页面 (ChatPage) | ⬜ 未开始 |
| C7 | 通讯录页面 | ⬜ 未开始 |

---

## 三、下一步

**当前应执行：Part A → Task A1（初始化 go-zero 项目骨架）**

1. 创建 `server/go.mod`、`server/Makefile`
2. 安装 go-zero + goctl
3. 创建各服务入口骨架

---

## 四、常用命令速查

```bash
# 中间件
cd /home/hcy/workspace/webchat/server
docker compose up -d        # 启动
docker compose down          # 停止
docker compose ps            # 查看状态

# Go
cd /home/hcy/workspace/webchat/server
go mod tidy                  # 整理依赖
make api                     # 生成 gateway 代码
make rpc                     # 生成 RPC 代码

# Flutter
cd /home/hcy/workspace/webchat/client
flutter pub get              # 安装依赖
flutter run                  # 运行
flutter build apk            # 打包
```

---

## 五、架构关键决策

- **微服务框架**：go-zero（非 gin），goctl 代码生成
- **认证**：Token + Redis（非 JWT），SHA256 生成，Redis 存储 24h TTL
- **对象存储**：MinIO
- **服务发现**：etcd（go-zero 内置）
- **WebSocket**：chat-service 独立 HTTP 端口 8081，不经过 gateway
- **文件结构**：go-zero 标准布局 `app/<service>/internal/{config,svc,dao,logic}`
