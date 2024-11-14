package main

import (
	"final-project/database"
	"final-project/handlers"
	"final-project/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         time.Duration
}

func getDefaultConfig() Config {
	return Config{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5674"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		MaxAge: 12 * time.Hour,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	db := database.ConnectDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	jwtKey := os.Getenv("JWT_KEY_SESSION")
	if jwtKey == "" {
		log.Fatal("JWT key environment variable is not set")
	}

	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     getDefaultConfig().AllowedOrigins,
		AllowMethods:     getDefaultConfig().AllowedMethods,
		AllowHeaders:     getDefaultConfig().AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           getDefaultConfig().MaxAge,
	}

	r.Use(cors.New(corsConfig))
	v1 := r.Group("/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"time":   time.Now().Format(time.RFC3339),
			})
		})

		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Welcome to the Deposito API",
				"version": "1.0",
			})
		})
		accountHandler := handlers.NewAccount(db, []byte(jwtKey))
		accountRoutes := r.Group("/account")
		{
			accountRoutes.POST("login/admin", accountHandler.AccountAdminLogin) // Restricted to admin login
			accountRoutes.POST("login/user", accountHandler.AccountUserLogin)
			accountRoutes.POST("signup", middleware.AuthJWTMiddleware(jwtKey), accountHandler.AccountSignUp)
			accountRoutes.POST("change-password", middleware.AuthJWTMiddleware(jwtKey), accountHandler.ChangePassword)
		}
		userHandler := handlers.NewUser(db)
		userRoutes := r.Group("/user")
		{
			userRoutes.GET("/profile", middleware.AuthJWTMiddleware(jwtKey), userHandler.Profile)
			userRoutes.GET("/mutation/transaction", middleware.AuthJWTMiddleware(jwtKey), userHandler.TransactionHistory)
			userRoutes.GET("/mutation/deposit", middleware.AuthJWTMiddleware(jwtKey), userHandler.PersonalDeposit)
			userRoutes.POST("/edit/profile", middleware.AuthJWTMiddleware(jwtKey), userHandler.EditProfile)
			userRoutes.POST("/register/deposit", middleware.AuthJWTMiddleware(jwtKey), userHandler.RegisterDeposit)
		}
		adminHandler := handlers.NewAdmin(db)
		adminRoutes := r.Group("/admin")
		{
			adminRoutes.GET("/list/user", adminHandler.ListUserProfile)
			adminRoutes.GET("/list/user/:id", adminHandler.DetailUser)
			adminRoutes.GET("/list/deposit/mutation", adminHandler.ListUserDeposito)
			adminRoutes.POST("/topup", adminHandler.TopUpUser)
		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
