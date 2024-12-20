package routes

import (
	"job_portal/api_gateway/internal/handler"
	// Import middleware package

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")

	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/verify", authHandler.VerifyUser)
		authRoutes.POST("/resend-otp", authHandler.ResendOtp)

		// Protected routes (require authMiddleware)
		authRoutes.GET("/user", authMiddleware, authHandler.GetUser)
		authRoutes.POST("/logout", authMiddleware, authHandler.Logout)
	}
}
