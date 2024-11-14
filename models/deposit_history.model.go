package model

import "time"

type DepositHistory struct {
	Id           int64     `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Deposit_Id   string    `json:"deposit_id"`
	Account_Id   int64     `json:"account_id"`
	Deposit_Name string    `json:"deposit_name"`
	Amount       int64     `json:"amount"`
	Time_Period  int       `json:"time_period"`
	Time_Stamp   time.Time `json:"time_stamp"`
}

func (DepositHistory) TableName() string {
	return "deposit_history"
}
