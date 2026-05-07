package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (*types.UserInfo, error) {
	uid := middleware.GetUserID(l.ctx)
	u, err := l.svcCtx.UserRpc.UpdateProfile(l.ctx, &userClient.UpdateProfileRequest{
		UserId: uid, Nickname: req.Nickname, Avatar: req.Avatar, Gender: int32(req.Gender), Region: req.Region, Signature: req.Signature,
	})
	if err != nil {
		return nil, err
	}
	return &types.UserInfo{Id: u.Id, Phone: u.Phone, Nickname: u.Nickname, Avatar: u.Avatar, Gender: int(u.Gender), Region: u.Region, Signature: u.Signature}, nil
}
