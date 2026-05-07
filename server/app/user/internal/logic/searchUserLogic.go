package logic

import (
	"context"
	"strconv"

	"github.com/team/webchat-server/app/user/internal/svc"
	"github.com/team/webchat-server/app/user/internal/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchUserLogic) SearchUser(in *user.SearchUserRequest) (*user.SearchUserResponse, error) {
	users, err := l.svcCtx.UserDao.Search(l.ctx, in.Keyword)
	if err != nil {
		return nil, err
	}
	var result []*user.UserInfo
	for _, u := range users {
		result = append(result, &user.UserInfo{
			Id: strconv.FormatInt(u.Id, 10), Phone: u.Phone, Nickname: u.Nickname,
			Avatar: u.Avatar, Gender: u.Gender, Region: u.Region, Signature: u.Signature, CreatedAt: u.CreatedAt,
		})
	}
	return &user.SearchUserResponse{Users: result}, nil
}
