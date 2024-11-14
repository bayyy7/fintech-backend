package model

type Account struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

func (Account) TableName() string {
	return "account"
}
