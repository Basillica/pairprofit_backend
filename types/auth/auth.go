package types

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     int64
}

type Todo struct {
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
}
