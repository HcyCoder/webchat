package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/team/webchat-server/common/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProfileLogic) GetProfile(in *user.GetProfileRequest) (*user.UserInfo, error) {
	id, _ := strconv.ParseInt(in.UserId, 10, 64)
	u, err := l.svcCtx.UserDao.FindById(l.ctx, id)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	return &user.UserInfo{
		Id: in.UserId, Phone: u.Phone, Nickname: u.Nickname,
		Avatar: u.Avatar, Gender: u.Gender, Region: u.Region, Signature: u.Signature, CreatedAt: u.CreatedAt,
	}, nil
}
