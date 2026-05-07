package dao

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/team/webchat-server/app/user/model"
)

type UserDao struct {
	conn sqlx.SqlConn
}

func NewUserDao(conn sqlx.SqlConn) *UserDao {
	return &UserDao{conn: conn}
}

func (d *UserDao) Insert(ctx context.Context, u *model.User) (int64, error) {
	u.CreatedAt = time.Now().UnixMilli()
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
