package dao

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ContactDao struct {
	conn sqlx.SqlConn
}

func NewContactDao(conn sqlx.SqlConn) *ContactDao {
	return &ContactDao{conn: conn}
}

type ContactRow struct {
	ContactId int64  `db:"contact_id"`
	Nickname  string `db:"nickname"`
	Avatar    string `db:"avatar"`
	Remark    string `db:"remark"`
	Tag       string `db:"tag"`
	IsBlocked bool   `db:"is_blocked"`
}

func (d *ContactDao) AddRequest(ctx context.Context, fromUser, toUser int64, message string) error {
	now := time.Now().UnixMilli()
	_, err := d.conn.ExecCtx(ctx,
		"INSERT INTO friend_requests (from_user, to_user, message, status, created_at) VALUES (?, ?, ?, 'pending', ?)",
		fromUser, toUser, message, now)
	return err
}

func (d *ContactDao) AcceptRequest(ctx context.Context, requestId int64) (fromUser, toUser int64, err error) {
	var req struct {
		Id       int64  `db:"id"`
		FromUser int64  `db:"from_user"`
		ToUser   int64  `db:"to_user"`
		Status   string `db:"status"`
	}
	err = d.conn.QueryRowCtx(ctx, &req,
		"SELECT id, from_user, to_user, status FROM friend_requests WHERE id = ? AND status = 'pending'", requestId)
	if err != nil {
		return 0, 0, err
	}
	_, err = d.conn.ExecCtx(ctx, "UPDATE friend_requests SET status = 'accepted' WHERE id = ?", requestId)
	if err != nil {
		return 0, 0, err
	}
	now := time.Now().UnixMilli()
	d.conn.ExecCtx(ctx, "INSERT IGNORE INTO contacts (user_id, contact_id, added_at) VALUES (?, ?, ?)", req.FromUser, req.ToUser, now)
	d.conn.ExecCtx(ctx, "INSERT IGNORE INTO contacts (user_id, contact_id, added_at) VALUES (?, ?, ?)", req.ToUser, req.FromUser, now)
	return req.FromUser, req.ToUser, nil
}

func (d *ContactDao) RejectRequest(ctx context.Context, requestId int64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE friend_requests SET status = 'rejected' WHERE id = ?", requestId)
	return err
}

func (d *ContactDao) List(ctx context.Context, userId int64) ([]*ContactRow, error) {
	var rows []*ContactRow
	err := d.conn.QueryRowsCtx(ctx, &rows,
		`SELECT c.contact_id, COALESCE(c.remark, u.nickname) AS nickname, u.avatar, COALESCE(c.remark, '') AS remark,
		        COALESCE(c.tag, '') AS tag, c.is_blocked
		 FROM contacts c JOIN users u ON c.contact_id = u.id WHERE c.user_id = ?`, userId)
	return rows, err
}
