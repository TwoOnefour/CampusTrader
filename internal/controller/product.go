package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/model"
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListProductSearchParams struct {
	LastId   uint64 `form:"last_id,default=0"`
	PageSize uint64 `form:"page_size,default=50"`
}

type ListProductSearchResult struct {
	List  []model.Product `json:"list"`
	Total uint64          `json:"total"`
	Page  uint64          `json:"page"`
	Size  uint64          `json:"size"`
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
}

func NewProductController(productSvc *service.ProductService) *ProductController {
	return &ProductController{
		productSvc: productSvc,
	}
}

// ListProducts GET /api/v1/product/list no login needed params: { "lastId": 123, "pageSize": 10 }
func (c *ProductController) ListProducts(ctx *gin.Context) {
	var req ListProductSearchParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	products, err := c.productSvc.ListProducts(ctx, req.PageSize, req.LastId)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, ListProductSearchResult{
		List:  products,
		Total: uint64(len(products)),
		Size:  req.PageSize,
	})
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

}
