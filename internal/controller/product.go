package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ListProductSearchResult = model.PageData[model.ProductWithUserRating]

type ListUserProductSearchResult = model.PageData[model.Product]

type SearchProductParams struct {
	Keyword string `form:"keyword"`
	Count   uint64 `form:"count"`
}

type DropProductReq struct {
	ProductId uint64 `form:"product_id"`
	Reason    string `form:"reason"`
}

type CreateProductReq struct {
	Name        string  `gorm:"column:name;type:VARCHAR(100);comment:商品名称;not null;" json:"name"`
	Description string  `gorm:"column:description;type:TEXT;comment:商品描述;not null;" json:"description"`
	Price       float64 `gorm:"column:price;type:DECIMAL(10, 2);comment:价格;not null;" json:"price"`
	CategoryId  uint64  `gorm:"column:category_id;type:BIGINT UNSIGNED;comment:分类ID;not null;" json:"category_id"`
	Status      string  `gorm:"column:status;type:ENUM('available', 'sold', 'removed');comment:状态;default:available;" json:"status"`
	Condition   string  `gorm:"column:condition;type:ENUM('new', 'like_new', 'good', 'fair', 'poor');comment:新旧程度;not null;" json:"condition"`
	ImageUrl    string  `gorm:"column:image_url;type:VARCHAR(255);comment:主图URL;" json:"image_url"`
}

type ProductController struct {
	productSvc *service.ProductService
	statsSvc   *service.StatisticsService
}

func NewProductController(productSvc *service.ProductService, statsSvc *service.StatisticsService) *ProductController {
	return &ProductController{
		productSvc: productSvc,
		statsSvc:   statsSvc,
	}
}

// ListProducts GET /api/v1/product/list no login needed params: { "lastId": 123, "pageSize": 10 }
func (c *ProductController) ListProducts(ctx *gin.Context) {
	var req model.PageParam
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	categoryId := ctx.Query("category")

	if categoryId != "" {
		var uintId uint64
		uintId, err := strconv.ParseUint(categoryId, 10, 64)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		products, err := c.productSvc.ListProductsByProc(ctx, uintId, req)
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(ctx, products)
		return
	}

	products, err := c.productSvc.ListProducts(ctx, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, products)
}

// CreateProduct POST /api/v1/product/create should login eq
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req CreateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := c.productSvc.CreateProduct(ctx, &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		Status:      req.Status,
		Condition:   req.Condition,
		ImageUrl:    req.ImageUrl,
		SellerId:    uint64(ctx.GetUint("userID")),
	})
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (c *ProductController) ListMyProducts(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	var req model.PageParam
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	products, err := c.productSvc.ListMyProducts(ctx, uint64(userID), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, products)
}

func (c *ProductController) SearchProducts(ctx *gin.Context) {
	var req model.PageParam
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	keyword := ctx.Query("keyword")
	products, err := c.productSvc.SearchProduct(ctx, req, keyword)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, products)
}

func (c *ProductController) SearchProductsSuggestion(ctx *gin.Context) {
	prefix := strings.TrimSpace(ctx.Query("prefix"))
	if prefix == "" {
		response.Error(ctx, http.StatusBadRequest, "不正确的参数")
		return
	}
	productPrefix, err := c.productSvc.SearchSuggestion(ctx, prefix)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, map[string]interface{}{
		"list": productPrefix,
	})
}

func (c *ProductController) DropProduct(ctx *gin.Context) {
	var dropProductReq DropProductReq
	if err := ctx.ShouldBindJSON(&dropProductReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.GetUint("userID")

	err := c.productSvc.DropProduct(ctx, dropProductReq.ProductId, uint64(userID), dropProductReq.Reason)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	response.Success(ctx, nil)
}
