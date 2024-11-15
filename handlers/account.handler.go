package handlers

import (
	model "final-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountInterface interface {
	AccountUserLogin(*gin.Context)
	AccountAdminLogin(*gin.Context)
	AccountSignUp(*gin.Context)
	ChangePassword(*gin.Context)
}

type accountImplement struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewAccount(db *gorm.DB, jwtKey []byte) AccountInterface {
	return &accountImplement{
		db,
		jwtKey,
	}
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *accountImplement) AccountAdminLogin(ctx *gin.Context) {
	payload := LoginPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account := model.Account{}
	if err := a.db.Where("username = ? AND role = ?", payload.Username, 1).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "username not found",
			})
			return
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "wrong password",
		})
		return
	}

	token, err := a.createJWT(&account)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}

func (a *accountImplement) AccountUserLogin(ctx *gin.Context) {
	payload := LoginPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account := model.Account{}
	if err := a.db.Where("username = ? AND role = ?", payload.Username, 0).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "username not found",
			})
			return
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "wrong password",
		})
		return
	}

	token, err := a.createJWT(&account)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}

type SignUpPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (a *accountImplement) AccountSignUp(ctx *gin.Context) {
	payload := SignUpPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingUser := model.Account{}
	if result := a.db.Where("username = ? AND role = ?", payload.Username, 0).First(&existingUser); result.RowsAffected > 0 {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "username already exist",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	newAccount := model.Account{
		Username: payload.Username,
		Password: string(hashPassword),
		Role:     0,
	}

	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&newAccount).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser := model.User{
		Account_Id: newAccount.Id,
		Name:       payload.Name,
	}

	if err := tx.Create(&newUser).Error; err != nil {
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

type accountChangePasswordPayload struct {
	Password string `json:"password"`
}

func (a *accountImplement) ChangePassword(ctx *gin.Context) {
	payload := accountChangePasswordPayload{}

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	id := ctx.GetInt64("id")
	var account model.Account
	if err := a.db.Where("role = ?", 0).First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "account not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account.Password = string(hashed)

	if err := a.db.Model(&account).Where("id = ?", id).Update("password", string(hashed)).Error; err != nil {
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
	})
}

func (a *accountImplement) createJWT(account *model.Account) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = account.Id
	claims["username"] = account.Username
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	tokenString, err := token.SignedString(a.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
