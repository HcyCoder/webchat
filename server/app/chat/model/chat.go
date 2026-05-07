package model

type Message struct {
	Id         int64  `db:"id"`
	ChatType   string `db:"chat_type"`
	FromUser   int64  `db:"from_user"`
	ToId       int64  `db:"to_id"`
	MsgType    string `db:"msg_type"`
	Content    string `db:"content"`
	IsRecalled bool   `db:"is_recalled"`
	CreatedAt  int64  `db:"created_at"`
}

type Conversation struct {
	Id           int64  `db:"id"`
	UserId       int64  `db:"user_id"`
	ChatType     string `db:"chat_type"`
	TargetId     int64  `db:"target_id"`
	LastMsgId    int64  `db:"last_msg_id"`
	UnreadCount  int32  `db:"unread_count"`
	IsPinned     bool   `db:"is_pinned"`
	IsMuted      bool   `db:"is_muted"`
	UpdatedAt    int64  `db:"updated_at"`
}

type ReadReceipt struct {
	MsgId  int64 `db:"msg_id"`
	UserId int64 `db:"user_id"`
	ReadAt int64 `db:"read_at"`
}
