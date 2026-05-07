package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefreshTokenLogic) RefreshToken() (*types.LoginResponse, error) {
	uid := middleware.GetUserID(l.ctx)
	tok, _ := l.svcCtx.TokenManager.Refresh(l.ctx, "", uid)
	return &types.LoginResponse{UserId: uid, Token: tok}, nil
}
