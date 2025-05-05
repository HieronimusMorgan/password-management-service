package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"password-management-service/internal/dto/in"
	"password-management-service/internal/services"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/jwt"
	"password-management-service/package/response"
)

type PasswordEntryController interface {
	AddPasswordEntry(context *gin.Context)
	GetPasswordEntryByID(context *gin.Context)
}

type passwordEntryController struct {
	PasswordEntryService services.PasswordEntryService
	JWTService           jwt.Service
}

func NewPasswordEntryController(passwordEntryService services.PasswordEntryService, jwtService jwt.Service) PasswordEntryController {
	return &passwordEntryController{
		PasswordEntryService: passwordEntryService,
		JWTService:           jwtService,
	}
}

func (c *passwordEntryController) AddPasswordEntry(context *gin.Context) {
	var req in.PasswordEntryRequest

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}
	//
	//credentialKey := context.GetHeader(utils.XCredentialKey)
	//if credentialKey == "" {
	//	response.SendResponse(context, http.StatusBadRequest, "Error", nil, "credential key not found")
	//	return
	//}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	if err := c.PasswordEntryService.AddPasswordEntry(&req, token.ClientID); err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", nil, "Password entry added successfully")
}

func (c *passwordEntryController) GetPasswordEntryByID(context *gin.Context) {
	entryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	passwordEntry, err := c.PasswordEntryService.GetPasswordEntryByID(entryID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", passwordEntry, nil)
}
