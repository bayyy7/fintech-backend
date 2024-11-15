package model

type User struct {
	Id             int64  `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Account_Id     int64  `json:"account_id"`
	Account_Number int64  `json:"account_number" gorm:"autoIncrement;<-:false"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Id_Card        int64  `json:"id_card"`
	Mothers_Name   string `json:"mothers_name"`
	Date_of_Birth  string `json:"date_of_birth"`
	Gender         string `json:"gender"`
	Balance        int64  `json:"balance"`
}

func (User) TableName() string {
	return "user"
}
