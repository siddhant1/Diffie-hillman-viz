package wrappers

type UserRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type User struct {
	Username       string `json:"username"`
	HashedPassword int    `json:"password"`
}
