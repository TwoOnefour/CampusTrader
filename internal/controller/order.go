package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderController struct {
	orderSvc *service.OrderService
}

type OrderReq struct {
	ItemId uint64 `json:"item_id" binding:"required"`
}

type ListOrderResult struct {
	List  []model.Order `json:"list"`
	Total uint64        `json:"total"`
}

func NewOrderController(svc *service.OrderService) *OrderController {
	return &OrderController{orderSvc: svc}
}

// need login
func (c *OrderController) Order(ctx *gin.Context) {
	var orderReq OrderReq
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	err := c.orderSvc.CreateOrder(ctx, orderReq.ItemId, ctx.GetUint64("userID"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
}

// need login /api/v1/order/my
func (c *OrderController) ListOrder(ctx *gin.Context) {
	userId := ctx.GetUint("userID")
	orders, err := c.orderSvc.ListOrder(ctx, userId)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, ListOrderResult{
		List:  orders,
		Total: uint64(len(orders)),
	})
}
