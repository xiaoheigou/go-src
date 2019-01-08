package response

type UserArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	//角色 0:管理员 1:坐席 2:平台商
	Role int `json:"role"`
}

type UserPasswordArgs struct {
	Username       string `json:"username"`
	OriginPassword string `json:"origin_password"`
	Password       string `json:"password"`
}
