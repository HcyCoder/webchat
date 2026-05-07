package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/team/webchat-server/app/user/model"
	"github.com/team/webchat-server/common/errcode"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
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
	u := &model.User{Phone: in.Phone, PasswordHash: string(hash), Nickname: in.Nickname}
	id, err := l.svcCtx.UserDao.Insert(l.ctx, u)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResponse{UserId: strconv.FormatInt(id, 10)}, nil
}
