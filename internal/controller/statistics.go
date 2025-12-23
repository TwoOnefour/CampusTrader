package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HotCategoriesResponse struct {
	List []model.Category `json:"list"`
}

type StatisticsController struct {
	svc *service.StatisticsService
}

func NewStatisticsController(svc *service.StatisticsService) *StatisticsController {
	return &StatisticsController{svc: svc}
}

func (s *StatisticsController) GetHotCategories(ctx *gin.Context) {
	categories, err := s.svc.GetHotCategories(ctx)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, HotCategoriesResponse{
		List: categories,
	})
}

func (s *StatisticsController) GetUserRating(ctx *gin.Context) {
	type UserParams struct {
		UserId uint64 `uri:"id" binding:"required,min=1"`
	}

	var params UserParams

	// 重点2：使用 ShouldBindUri 解析路径参数
	if err := ctx.ShouldBindUri(&params); err != nil {
		// 如果传进来的不是数字，或者 <=0，这里直接报错
		response.Error(ctx, http.StatusBadRequest, "Invalid User ID")
		return
	}

	// 调用 Service
	rating, err := s.svc.GetUserRating(ctx, params.UserId)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, map[string]interface{}{
		"rating": rating,
	})
}

func (s *StatisticsController) GetUserCompletedOrderRecord(ctx *gin.Context) {
	type userParams struct {
		UserId uint64 `uri:"id" binding:"required,min=1"`
	}
	var userReq userParams
	if err := ctx.ShouldBindUri(&userReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	orders, err := s.svc.GetUserCompletedOrderRecord(ctx, userReq.UserId)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	type res struct {
		List []model.Order `json:"list"`
	}
	response.Success(ctx, &res{
		List: orders,
	})
}
