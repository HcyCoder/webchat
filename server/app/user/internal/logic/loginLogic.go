package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/team/webchat-server/common/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
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
