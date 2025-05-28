package routes

import (
	"github.com/gin-gonic/gin"
	"password-management-service/config"
	"password-management-service/internal/controller"
)

func PasswordEntryRoutes(r *gin.Engine, middleware config.Middleware, controller controller.PasswordEntryController) {
	routerGroup := r.Group("/v1/entry")
	routerGroup.Use(middleware.PasswordMiddleware.HandlerPassword())
	{
		routerGroup.POST("/", controller.AddPasswordEntry)
		routerGroup.PUT("/:id", controller.UpdatePasswordEntry)
		routerGroup.POST("/group/:id", controller.AddGroupPasswordEntry)
		routerGroup.GET("/", controller.GetListPasswordEntries)
		routerGroup.GET("/:id", controller.GetPasswordEntryByID)
		routerGroup.DELETE("/:id", controller.DeletePasswordEntry)
	}
}
