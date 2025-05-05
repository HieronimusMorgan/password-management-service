package routes

import (
	"github.com/gin-gonic/gin"
	"password-management-service/config"
	"password-management-service/internal/controller"
)

func PasswordRoutes(r *gin.Engine, middleware config.Middleware, controller controller.PasswordEntryController) {

	routerGroup := r.Group("/v1/password")
	routerGroup.Use(middleware.PasswordMiddleware.HandlerPassword())
	{
		routerGroup.POST("/add", controller.AddPasswordEntry)
		routerGroup.GET("/:id", controller.GetPasswordEntryByID)
	}
}
