package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"password-management-service/internal/services"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/jwt"
	"password-management-service/package/response"
)

type PasswordTagController interface {
	AddPasswordTag(context *gin.Context)
	UpdatePasswordTag(context *gin.Context)
	GetListPasswordTag(context *gin.Context)
	DeletePasswordTag(context *gin.Context)
}

type passwordTagController struct {
	PasswordTagService services.PasswordTagService
	JWT                jwt.Service
}

func NewPasswordTagController(passwordTagService services.PasswordTagService, JWT jwt.Service) PasswordTagController {
	return &passwordTagController{
		PasswordTagService: passwordTagService,
		JWT:                JWT,
	}
}

func (p *passwordTagController) AddPasswordTag(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	passwordTag, err := p.PasswordTagService.AddPasswordTag(req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Failed to add password tag", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusCreated, "Password tag added successfully", passwordTag, nil)
}

func (p *passwordTagController) UpdatePasswordTag(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	tagID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	passwordTag, err := p.PasswordTagService.UpdatePasswordTag(tagID, req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Failed to update password tag", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Password tag updated successfully", passwordTag, nil)
}

func (p *passwordTagController) GetListPasswordTag(context *gin.Context) {
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

	tags, total, err := p.PasswordTagService.GetListPasswordTag(token.ClientID, pageIndex, pageSize)
	if err != nil {
		response.SendResponseList(context, 500, "Failed to get list password tag", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}

	response.SendResponseList(context, 200, "Get list password tag successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     tags,
	}, nil)
}

func (p *passwordTagController) DeletePasswordTag(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	tagID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	err = p.PasswordTagService.DeletePasswordTagByID(tagID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Failed to delete password tag", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Password tag deleted successfully", nil, nil)
}
