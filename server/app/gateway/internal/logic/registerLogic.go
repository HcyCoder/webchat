package logic

import (
	"context"

	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/user/userClient"
	"github.com/team/webchat-server/common/token"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.RegisterResponse, error) {
	resp, err := l.svcCtx.UserRpc.Register(l.ctx, &userClient.RegisterRequest{
		Phone: req.Phone, Password: req.Password, Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}
	tok := token.Generate(resp.UserId)
	if err := l.svcCtx.TokenManager.Store(l.ctx, tok, resp.UserId); err != nil {
		return nil, err
	}
	return &types.RegisterResponse{UserId: resp.UserId, Token: tok}, nil
}
