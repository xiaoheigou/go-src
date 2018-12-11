package models

import "yuudidi.com/pkg/utils"

type User struct {
	ID
	Username  string `gorm:"unique;type:varchar(50)" json:"username"`
	Password  []byte `gorm:"type:varbinary(64);column:password;not null" json:"-"`
	Salt      []byte `gorm:"type:varbinary(64);column:salt" json:"-"`
	Algorithm string `gorm:"type:varchar(255)" json:"-"`
	Phone     string `gorm:"type:varchar(30)" json:"phone"`
	Email     string `gorm:"type:varchar(50)" json:"email"`
	Address   string `gorm:"type:varchar(200)" json:"address"`
	//用户状态 0: 正常 1: 冻结
	UserStatus int `gorm:"type:tinyint(1);default:0" json:"user_status"`
	//平台角色 0:管理员 1:坐席
	Role int `gorm:"type:tinyint(1)" json:"role"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&User{})
}
