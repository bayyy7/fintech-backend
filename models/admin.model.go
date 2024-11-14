package model

type Admin struct {
	Id         int64  `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Account_Id int64  `json:"account_id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
}

func (Admin) TableName() string {
	return "admin"
}
