package handlers

import (
	"bytes"
	"encoding/json"
	model "final-project/models"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInterface interface {
	Profile(*gin.Context)
	TransactionHistory(*gin.Context)
	DepositHistory(*gin.Context)
	EditProfile(*gin.Context)
	RegisterDeposit(*gin.Context)
	PersonalDeposit(*gin.Context)
}

type userImplement struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UserInterface {
	return &userImplement{
		db,
	}
}

func (a *userImplement) Profile(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var user model.User

	if err := a.db.Where("role = ? AND account_id = ?", 0, id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// PERLU DIBERI RENTANG TANGGAL
func (a *userImplement) TransactionHistory(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var mutation []model.TransactionHistory

	if err := a.db.Find(&mutation).Where("account_id = ?", id).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": mutation,
	})
}

func (a *userImplement) DepositHistory(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var mutation []model.DepositHistory

	if err := a.db.Find(&mutation).Where("account_id = ?", id).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": mutation,
	})
}

type ProfilePayload struct {
	Address       string    `json:"address"`
	Id_Card       int64     `json:"id_card"`
	Mothers_Name  string    `json:"mothers_name"`
	Date_of_Birth time.Time `json:"date_of_birth"`
	Gender        string    `json:"gender"`
	Balance       int64     `json:"balance"`
}

func (a *userImplement) EditProfile(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	payload := ProfilePayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user model.User
	if err := a.db.Where("role = ? AND account_id = ?", 0, id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := a.db.Model(&user).Where("account_id = ?", id).Updates(payload).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})
}

type DepositPayload struct {
	Deposit_Id string `json:"deposit_id" binding:"required"`
	Account_Id int64  `json:"account_id"`
	Name       string `json:"name"`
	Amount     int64  `json:"amount"`
	Min_Month  int    `json:"min_amount"`
}

func (a *userImplement) RegisterDeposit(ctx *gin.Context) {
	payload := DepositPayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch payload.Deposit_Id {
	case "mini":
		if payload.Amount < 100000 && payload.Amount > 10000000 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "amount not acceptable",
			})
			return
		}
	case "maxi":
		if payload.Amount < 100000 && payload.Amount > 1000000000 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "amount not acceptable",
			})
			return
		}
	case "great":
		if payload.Amount < 1000000000 && payload.Amount > 9999999999 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "amount not acceptable",
			})
			return
		}
	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "deposit type not recognized",
		})
		return
	}

	id := ctx.GetInt64("id")
	var user model.User
	if err := a.db.Where("role = ? AND account_id = ?", 0, id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	server := os.Getenv("SERVER_API")
	apiURL := server + "/deposito"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to marshal JSON",
		})
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create request",
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "API request failed",
			"status": resp.StatusCode,
		})
		return
	}

	currentDate, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newDepositHistory := model.DepositHistory{
		Deposit_Id:   payload.Deposit_Id,
		Account_Id:   id,
		Deposit_Name: payload.Name,
		Amount:       payload.Amount,
		Time_Period:  payload.Min_Month,
		Time_Stamp:   currentDate,
	}

	newTransactionHistory := model.TransactionHistory{
		Account_Id:           id,
		Transaction_Category: "Deposito",
		Amount:               payload.Amount,
		In_Out:               1,
		Time_Stamp:           time.Now(),
	}

	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&newDepositHistory).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := tx.Create(&newTransactionHistory).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (a *userImplement) PersonalDeposit(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var user model.User
	if err := a.db.Where("role = ? AND account_id = ?", 0, id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	server := os.Getenv("SERVER_API") + "/deposito/" + strconv.FormatInt(id, 10)
	response, err := http.Get(server)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create request",
		})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "API request failed",
			"status": response.StatusCode,
		})
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data struct {
		Data []model.Deposit `json:"data"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data.Data,
	})
}
