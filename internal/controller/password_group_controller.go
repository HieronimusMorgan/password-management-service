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

type PasswordGroupController interface {
	AddPasswordGroup(context *gin.Context)
	UpdatePasswordGroup(context *gin.Context)
	GetListPasswordGroup(context *gin.Context)
	GetItemListPasswordGroup(context *gin.Context)
	DeletePasswordGroup(context *gin.Context)
}

type passwordGroupController struct {
	PasswordGroupService services.PasswordGroupService
	JWT                  jwt.Service
}

func NewPasswordGroupController(passwordGroupService services.PasswordGroupService, JWT jwt.Service) PasswordGroupController {
	return &passwordGroupController{
		PasswordGroupService: passwordGroupService,
		JWT:                  JWT,
	}
}

func (c *passwordGroupController) AddPasswordGroup(context *gin.Context) {
	var req in.PasswordGroupRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	passwordGroup, err := c.PasswordGroupService.AddPasswordGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", passwordGroup, nil)
}

func (c *passwordGroupController) UpdatePasswordGroup(context *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	passwordGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	passwordGroup, err := c.PasswordGroupService.UpdatePasswordGroup(passwordGroupID, req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", passwordGroup, nil)
}

func (c *passwordGroupController) GetListPasswordGroup(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	passwordGroups, err := c.PasswordGroupService.GetListPasswordGroup(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", passwordGroups, nil)
}

func (c *passwordGroupController) GetItemListPasswordGroup(context *gin.Context) {
	passwordGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	passwordGroup, err := c.PasswordGroupService.GetItemListPasswordGroup(passwordGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", passwordGroup, nil)
}

func (c *passwordGroupController) DeletePasswordGroup(context *gin.Context) {
	passwordGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = c.PasswordGroupService.DeletePasswordGroupByID(passwordGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, nil)
}
