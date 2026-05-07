package logic

import (
	"context"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/media/mediaClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadURLLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUploadURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadURLLogic {
	return &GetUploadURLLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUploadURLLogic) GetUploadURL() (*types.UploadURLResponse, error) {
	resp, err := l.svcCtx.MediaRpc.GetUploadURL(l.ctx, &mediaClient.GetUploadURLRequest{})
	if err != nil {
		return nil, err
	}
	return &types.UploadURLResponse{UploadUrl: resp.UploadUrl, FileId: resp.FileId, ExpiresIn: resp.ExpiresIn}, nil
}
