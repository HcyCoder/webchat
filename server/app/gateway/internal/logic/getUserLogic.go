package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserLogic) GetUser() (*types.UserInfo, error) {
	uid := middleware.GetUserID(l.ctx)
	u, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &userClient.GetProfileRequest{UserId: uid})
	if err != nil {
		return nil, err
	}
	return &types.UserInfo{Id: u.Id, Phone: u.Phone, Nickname: u.Nickname, Avatar: u.Avatar, Gender: int(u.Gender), Region: u.Region, Signature: u.Signature}, nil
}
