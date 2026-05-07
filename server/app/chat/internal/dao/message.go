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
		"SELECT id, chat_type, from_user, to_id, msg_type, content, is_recalled, created_at FROM messages WHERE to_id=? AND chat_type=? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		toId, chatType, pageSize, offset)
	return msgs, err
}

func (d *MessageDao) Recall(ctx context.Context, msgId, userId int64) error {
	_, err := d.conn.ExecCtx(ctx,
		"UPDATE messages SET is_recalled=1 WHERE id=? AND from_user=?", msgId, userId)
	return err
}
