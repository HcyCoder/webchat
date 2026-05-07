package logic

import (
	"context"

	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/team/webchat-server/common/token"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	resp, err := l.svcCtx.UserRpc.Login(l.ctx, &userClient.LoginRequest{
		Phone: req.Phone, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	tok := token.Generate(resp.UserId)
	if err := l.svcCtx.TokenManager.Store(l.ctx, tok, resp.UserId); err != nil {
		return nil, err
	}
	return &types.LoginResponse{UserId: resp.UserId, Token: tok}, nil
}
