package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"password-management-service/internal/utils/jwt"
	"password-management-service/package/response"
)

type AdminMiddleware interface {
	HandlerAsset() gin.HandlerFunc
}

type adminMiddleware struct {
	JWTService jwt.Service
}

func NewAdminMiddleware(jwtService jwt.Service) AdminMiddleware {
	return adminMiddleware{
		JWTService: jwtService,
	}
}

func (a adminMiddleware) HandlerAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateTokenAdmin(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		tokenClaims, err := a.JWTService.ExtractClaims(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token claims", nil, err.Error())
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)

		c.Next()
	}
}
