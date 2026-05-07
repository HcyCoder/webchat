package dao

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/team/webchat-server/app/group/model"
)

type GroupDao struct{ conn sqlx.SqlConn }

func NewGroupDao(conn sqlx.SqlConn) *GroupDao { return &GroupDao{conn} }

func (d *GroupDao) Create(ctx context.Context, g *model.Group) (int64, error) {
	now := time.Now().UnixMilli()
	result, err := d.conn.ExecCtx(ctx,
		"INSERT INTO grps (name, avatar, owner_id, member_count, max_members, created_at) VALUES (?, ?, ?, 1, 500, ?)",
		g.Name, g.Avatar, g.OwnerId, now)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	d.conn.ExecCtx(ctx, "INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'owner', ?)", id, g.OwnerId, now)
	return id, nil
}

func (d *GroupDao) FindById(ctx context.Context, id int64) (*model.Group, error) {
	var g model.Group
	err := d.conn.QueryRowCtx(ctx, &g,
		"SELECT id, name, avatar, owner_id, announcement, member_count, max_members, created_at FROM grps WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (d *GroupDao) AddMembers(ctx context.Context, groupId int64, userIds []int64) error {
	now := time.Now().UnixMilli()
	for _, uid := range userIds {
		d.conn.ExecCtx(ctx, "INSERT IGNORE INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'member', ?)", groupId, uid, now)
	}
	d.conn.ExecCtx(ctx, "UPDATE grps SET member_count=(SELECT COUNT(*) FROM group_members WHERE group_id=?) WHERE id=?", groupId, groupId)
	return nil
}

func (d *GroupDao) RemoveMember(ctx context.Context, groupId, userId int64) error {
	d.conn.ExecCtx(ctx, "DELETE FROM group_members WHERE group_id=? AND user_id=? AND role!='owner'", groupId, userId)
	d.conn.ExecCtx(ctx, "UPDATE grps SET member_count=(SELECT COUNT(*) FROM group_members WHERE group_id=?) WHERE id=?", groupId, groupId)
	return nil
}

type MemberRow struct {
	UserId   int64  `db:"user_id"`
	Role     string `db:"role"`
	Alias    string `db:"alias"`
	IsMuted  bool   `db:"is_muted"`
	JoinedAt int64  `db:"joined_at"`
}

func (d *GroupDao) ListMembers(ctx context.Context, groupId int64) ([]*MemberRow, error) {
	var rows []*MemberRow
	err := d.conn.QueryRowsCtx(ctx, &rows,
		"SELECT user_id, role, alias, is_muted, joined_at FROM group_members WHERE group_id=?", groupId)
	return rows, err
}
