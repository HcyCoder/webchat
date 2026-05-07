package dao

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConversationDao struct{ conn sqlx.SqlConn }

func NewConversationDao(conn sqlx.SqlConn) *ConversationDao { return &ConversationDao{conn} }

func (d *ConversationDao) Upsert(ctx context.Context, userId, targetId int64, chatType string, lastMsgId int64) error {
	now := time.Now().UnixMilli()
	_, err := d.conn.ExecCtx(ctx,
		`INSERT INTO conversations (user_id, chat_type, target_id, last_msg_id, unread_count, updated_at)
		 VALUES (?, ?, ?, ?, 1, ?)
		 ON DUPLICATE KEY UPDATE last_msg_id=?, unread_count=unread_count+1, updated_at=?`,
		userId, chatType, targetId, lastMsgId, now, lastMsgId, now)
	return err
}

func (d *ConversationDao) ListByUser(ctx context.Context, userId int64) ([]*struct {
	Id         int64  `db:"id"`
	ChatType   string `db:"chat_type"`
	TargetId   int64  `db:"target_id"`
	LastMsgId  int64  `db:"last_msg_id"`
	UnreadCnt  int32  `db:"unread_count"`
	IsPinned   bool   `db:"is_pinned"`
	IsMuted    bool   `db:"is_muted"`
	UpdatedAt  int64  `db:"updated_at"`
}, error) {
	var rows []*struct {
		Id         int64  `db:"id"`
		ChatType   string `db:"chat_type"`
		TargetId   int64  `db:"target_id"`
		LastMsgId  int64  `db:"last_msg_id"`
		UnreadCnt  int32  `db:"unread_count"`
		IsPinned   bool   `db:"is_pinned"`
		IsMuted    bool   `db:"is_muted"`
		UpdatedAt  int64  `db:"updated_at"`
	}
	err := d.conn.QueryRowsCtx(ctx, &rows,
		"SELECT id, chat_type, target_id, last_msg_id, unread_count, is_pinned, is_muted, updated_at FROM conversations WHERE user_id=? ORDER BY updated_at DESC", userId)
	return rows, err
}

func (d *ConversationDao) ClearUnread(ctx context.Context, convId, userId int64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE conversations SET unread_count=0 WHERE id=? AND user_id=?", convId, userId)
	return err
}
