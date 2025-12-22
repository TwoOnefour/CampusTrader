package controller

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageController struct {
	svc *service.ImageService
}

func NewImageController(svc *service.ImageService) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

// upload need login
func (c *ImageController) Upload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 1. 打开文件流 (得到 multipart.File，它实现了 io.Reader)
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// 2. 生成路径
	ext := filepath.Ext(fileHeader.Filename)
	fileName := uuid.New().String() + ext

	// 3. 调用接口上传 (不用管是存本地还是 S3)
	url, err := c.svc.Save(ctx, file, fileName, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		response.Error(ctx, 500, "上传失败")
		return
	}

	response.Success(ctx, map[string]string{"url": url})
}
