package model

import "time"

type TransactionHistory struct {
	Id                   int64     `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Account_Id           int64     `json:"account_id"`
	Transaction_Category string    `json:"transaction_category"`
	Amount               int64     `json:"amount"`
	In_Out               int       `json:"in_out"`
	Time_Stamp           time.Time `json:"time_stamp"`
}

func (TransactionHistory) TableName() string {
	return "transaction_history"
}
