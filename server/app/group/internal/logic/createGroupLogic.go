package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/group/internal/group"
	"github.com/team/webchat-server/app/group/internal/svc"
	"github.com/team/webchat-server/app/group/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateGroupLogic) CreateGroup(in *group.CreateGroupRequest) (*group.GroupInfo, error) {
	ownerId, _ := strconv.ParseInt(in.OwnerId, 10, 64)
	g := &model.Group{Name: in.Name, OwnerId: ownerId}
	id, err := l.svcCtx.GroupDao.Create(l.ctx, g)
	if err != nil {
		return nil, err
	}
	var mids []int64
	for _, s := range in.MemberIds {
		mid, _ := strconv.ParseInt(s, 10, 64)
		mids = append(mids, mid)
	}
	l.svcCtx.GroupDao.AddMembers(l.ctx, id, mids)
	return &group.GroupInfo{
		Id: strconv.FormatInt(id, 10), Name: in.Name, OwnerId: in.OwnerId,
		MemberCount: 1 + int32(len(mids)), MaxMembers: 500, CreatedAt: g.CreatedAt,
	}, nil
}
