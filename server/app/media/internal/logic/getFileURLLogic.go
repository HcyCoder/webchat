package logic

import (
	"context"
	"time"

	"github.com/team/webchat-server/app/media/internal/media"
	"github.com/team/webchat-server/app/media/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileURLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileURLLogic {
	return &GetFileURLLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFileURLLogic) GetFileURL(in *media.GetFileURLRequest) (*media.GetFileURLResponse, error) {
	url, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, l.svcCtx.Config.Minio.Bucket, in.FileId, 1*time.Hour, nil)
	if err != nil {
		return nil, err
	}
	return &media.GetFileURLResponse{Url: url.String()}, nil
}
