package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"password-management-service/internal/utils/jwt"
	"password-management-service/package/response"
)

type PasswordMiddleware interface {
	HandlerPassword() gin.HandlerFunc
}

type passwordMiddleware struct {
	JWTService jwt.Service
}

func NewPasswordMiddleware(jwtService jwt.Service) PasswordMiddleware {
	return passwordMiddleware{
		JWTService: jwtService,
	}
}

func (a passwordMiddleware) HandlerPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateToken(token)
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

		if tokenClaims.Authorized == false {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		if tokenClaims.Exp < jwt.GetCurrentTime() {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "Token has expired")
			c.Abort()
			return
		}

		if !jwt.HasPasswordResource(tokenClaims.Resource) {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)
		c.Next()
	}
}
