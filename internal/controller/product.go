package controller

import (
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productSvc *service.ProductService
}

func NewProductController(productSvc *service.ProductService) *ProductController {
	return &ProductController{
		productSvc: productSvc,
	}
}

// ListProducts GET /api/v1/product/list no login needed
func (c *ProductController) ListProducts(ctx *gin.Context) {
	return c.productSvc
}

// CreateProduct POST /api/v1/product/create should login
func (c *ProductController) CreateProduct(ctx *gin.Context) {

}
