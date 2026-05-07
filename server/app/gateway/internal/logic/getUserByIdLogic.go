package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserByIdLogic) GetUserById(req *types.UserInfo) (*types.UserInfo, error) {
	u, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userClient.GetUserByIdRequest{UserId: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.UserInfo{Id: u.Id, Phone: u.Phone, Nickname: u.Nickname, Avatar: u.Avatar, Gender: int(u.Gender), Region: u.Region, Signature: u.Signature}, nil
}
