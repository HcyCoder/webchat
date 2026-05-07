package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/team/webchat-server/app/media/internal/media"
	"github.com/team/webchat-server/app/media/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadURLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadURLLogic {
	return &GetUploadURLLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetUploadURLLogic) GetUploadURL(in *media.GetUploadURLRequest) (*media.GetUploadURLResponse, error) {
	fileID := uuid.New().String()
	objectName := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), fileID)
	url, err := l.svcCtx.MinioClient.PresignedPutObject(l.ctx, l.svcCtx.Config.Minio.Bucket, objectName, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	return &media.GetUploadURLResponse{UploadUrl: url.String(), FileId: fileID, ExpiresIn: 900}, nil
}
