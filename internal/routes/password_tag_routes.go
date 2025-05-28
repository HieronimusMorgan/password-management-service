package routes

import (
	"github.com/gin-gonic/gin"
	"password-management-service/config"
	"password-management-service/internal/controller"
)

func PasswordTagRoutes(r *gin.Engine, middleware config.Middleware, controller controller.PasswordTagController) {
	routerTag := r.Group("/v1/tag")
	routerTag.Use(middleware.PasswordMiddleware.HandlerPassword())
	{
		routerTag.POST("/", controller.AddPasswordTag)
		routerTag.PUT("/:id", controller.UpdatePasswordTag)
		routerTag.GET("/", controller.GetListPasswordTag)
		routerTag.DELETE("/:id", controller.DeletePasswordTag)
	}
}
