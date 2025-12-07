package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listCategoryResult struct {
	List  []model.Category `json:"list"`
	Total uint64           `json:"total"`
}

type CategoryController struct {
	svc *service.CategoryService
}

func NewCategoryController(svc *service.CategoryService) *CategoryController {
	return &CategoryController{
		svc: svc,
	}
}

func (c *CategoryController) ListCategory(ctx *gin.Context) {
	category, err := c.svc.ListCategory(ctx)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, listCategoryResult{
		List:  category,
		Total: uint64(len(category)),
	})
}
