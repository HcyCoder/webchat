package model

type Group struct {
	Id           int64  `db:"id"`
	Name         string `db:"name"`
	Avatar       string `db:"avatar"`
	OwnerId      int64  `db:"owner_id"`
	Announcement string `db:"announcement"`
	MemberCount  int32  `db:"member_count"`
	MaxMembers   int32  `db:"max_members"`
	CreatedAt    int64  `db:"created_at"`
}

type GroupMember struct {
	GroupId  int64  `db:"group_id"`
	UserId   int64  `db:"user_id"`
	Role     string `db:"role"`
	Alias    string `db:"alias"`
	IsMuted  bool   `db:"is_muted"`
	JoinedAt int64  `db:"joined_at"`
}
