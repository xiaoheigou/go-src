package response

type UserArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Role     int    `json:"-"`
}

type UserPasswordArgs struct {
	Username       string `json:"username"`
	OriginPassword string `json:"origin_password"`
	Password       string `json:"password"`
}
