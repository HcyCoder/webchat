package model

type User struct {
	Id           int64  `db:"id"`
	Phone        string `db:"phone"`
	PasswordHash string `db:"password_hash"`
	Nickname     string `db:"nickname"`
	Avatar       string `db:"avatar"`
	Gender       int32  `db:"gender"`
	Region       string `db:"region"`
	Signature    string `db:"signature"`
	CreatedAt    int64  `db:"created_at"`
}
