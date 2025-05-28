package routes

import (
	"github.com/gin-gonic/gin"
	"password-management-service/config"
	"password-management-service/internal/controller"
)

func PasswordGroupRoutes(r *gin.Engine, middleware config.Middleware, controller controller.PasswordGroupController) {
	routerGroup := r.Group("/v1/group")
	routerGroup.Use(middleware.PasswordMiddleware.HandlerPassword())
	{
		routerGroup.POST("/", controller.AddPasswordGroup)
		routerGroup.PUT("/:id", controller.UpdatePasswordGroup)
		routerGroup.GET("/", controller.GetListPasswordGroup)
		routerGroup.GET("/item/:id", controller.GetItemListPasswordGroup)
		routerGroup.DELETE("/:id", controller.DeletePasswordGroup)
	}
}
