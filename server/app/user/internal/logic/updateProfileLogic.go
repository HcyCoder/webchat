package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/team/webchat-server/app/user/model"
	"github.com/team/webchat-server/common/errcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateProfileLogic) UpdateProfile(in *user.UpdateProfileRequest) (*user.UserInfo, error) {
	id, _ := strconv.ParseInt(in.UserId, 10, 64)
	u := &model.User{
		Id: id, Nickname: in.Nickname, Avatar: in.Avatar,
		Gender: in.Gender, Region: in.Region, Signature: in.Signature,
	}
	if err := l.svcCtx.UserDao.Update(l.ctx, u); err != nil {
		return nil, errcode.ErrUserNotFound
	}
	updated, _ := l.svcCtx.UserDao.FindById(l.ctx, id)
	return &user.UserInfo{
		Id: in.UserId, Phone: updated.Phone, Nickname: updated.Nickname,
		Avatar: updated.Avatar, Gender: updated.Gender, Region: updated.Region,
		Signature: updated.Signature, CreatedAt: updated.CreatedAt,
	}, nil
}
