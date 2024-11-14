package handlers

import (
	model "final-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminInterface interface {
	ListUserProfile(*gin.Context)
	DetailUser(*gin.Context)
	ListUserDeposito(*gin.Context)
	TopUpUser(*gin.Context)
}

type adminImplement struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) AdminInterface {
	return &adminImplement{
		db,
	}
}

func (a *adminImplement) ListUserProfile(ctx *gin.Context) {
	var user []model.User

	if err := a.db.Find(&user).Where("role = ?", 0).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (a *adminImplement) DetailUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user model.User

	if err := a.db.First(&user, id).Where("role = ?", 0).Error; err != nil {
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

func (a *adminImplement) ListUserDeposito(ctx *gin.Context) {
	var deposit_history []model.DepositHistory

	if err := a.db.Find(&deposit_history).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": deposit_history,
	})
}

type TransferPayload struct {
	Username string `json:"username" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
}

func (a *adminImplement) TopUpUser(ctx *gin.Context) {
	payload := TransferPayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	account := model.Account{}
	if result := a.db.Where("username = ? AND role = ?", payload.Username, 0).First(&account); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		return
	}

	user := model.User{}
	if result := a.db.Where("account_id = ?", account.Id).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		return
	}

	newTransactionHistory := model.TransactionHistory{
		Account_Id:           user.Account_Id,
		Transaction_Category: "TopUp",
		Amount:               payload.Amount,
		In_Out:               0,
		Time_Stamp:           time.Now(),
	}

	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.Balance += payload.Amount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := tx.Save(&newTransactionHistory).Error; err != nil {
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
		"balance": user.Balance,
	})
}
