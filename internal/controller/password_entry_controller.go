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
	UpdatePasswordEntry(context *gin.Context)
	AddGroupPasswordEntry(context *gin.Context)
	GetListPasswordEntries(context *gin.Context)
	GetPasswordEntryByID(context *gin.Context)
	DeletePasswordEntry(context *gin.Context)
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

	requestID := context.GetHeader(utils.XRequestID)
	if requestID == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "request id not found")
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	if err := c.PasswordEntryService.AddPasswordEntry(&req, token.ClientID, requestID); err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", nil, "Password entry added successfully")
}

func (c *passwordEntryController) UpdatePasswordEntry(context *gin.Context) {
	entryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	var req in.PasswordEntryRequest

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	if err := c.PasswordEntryService.UpdatePasswordEntry(entryID, &req, token.ClientID); err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", nil, "Password entry updated successfully")
}

func (c *passwordEntryController) AddGroupPasswordEntry(context *gin.Context) {
	var req struct {
		GroupID uint `json:"group_id"`
		EntryID uint `json:"entry_id"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	if err := c.PasswordEntryService.AddGroupPasswordEntry(req, token.ClientID); err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", nil, "Password entry updated successfully")
}

func (c *passwordEntryController) GetListPasswordEntries(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	pageIndex, pageSize, err := utils.GetPageIndexPageSize(context)
	if err != nil {
		response.SendResponse(context, 400, "Invalid page index or page size", nil, err.Error())
		return
	}

	tagsParam := context.Query("tags")

	passwordEntries, total, err := c.PasswordEntryService.GetListPasswordEntries(token.ClientID, tagsParam, pageIndex, pageSize)
	if err != nil {
		response.SendResponseList(context, 500, "Failed to get list password entry", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}

	response.SendResponseList(context, 200, "Get list password entry successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     passwordEntries,
	}, nil)
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

func (c *passwordEntryController) DeletePasswordEntry(context *gin.Context) {
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

	if err := c.PasswordEntryService.DeletePasswordEntry(entryID, token.ClientID); err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", nil, "Password entry deleted successfully")
}
