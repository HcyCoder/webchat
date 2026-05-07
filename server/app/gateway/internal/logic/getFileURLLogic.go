package logic

import (
	"context"
	"net/http"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/team/webchat-server/app/gateway/internal/types"
	"github.com/team/webchat-server/app/media/mediaClient"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileURLLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileURLLogic {
	return &GetFileURLLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFileURLLogic) GetFileURL(r *http.Request) (*types.FileURLResponse, error) {
	fid := extractPathParam(r.URL.Path, "files/")
	resp, err := l.svcCtx.MediaRpc.GetFileURL(l.ctx, &mediaClient.GetFileURLRequest{FileId: fid})
	if err != nil {
		return nil, err
	}
	return &types.FileURLResponse{Url: resp.Url}, nil
}

