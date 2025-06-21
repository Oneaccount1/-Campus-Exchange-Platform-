package controllers

import (
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"campus/internal/utils/upload"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

// UploadController 上传控制器
type UploadController struct{}

// NewUploadController 创建上传控制器
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage 上传图片
func (c *UploadController) UploadImage(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("获取上传文件失败", err))
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		response.HandleError(ctx, errors.NewBadRequestError("只支持jpg/jpeg/png/gif格式的图片", nil))
		return
	}

	// 保存文件
	fileType := "images" // 文件类型子目录
	filePath, err := upload.SaveUploadedFile(file, fileType)
	if err != nil {
		response.HandleError(ctx, errors.NewInternalServerError("保存文件失败", err))
		return
	}

	// 返回文件路径
	response.SuccessWithMessage(ctx, "上传成功", gin.H{
		"url": filePath,
	})
}

// UploadFile 上传通用文件
func (c *UploadController) UploadFile(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("获取上传文件失败", err))
		return
	}

	// 检查文件类型参数
	fileType := ctx.DefaultQuery("type", "files")
	
	// 保存文件
	filePath, err := upload.SaveUploadedFile(file, fileType)
	if err != nil {
		response.HandleError(ctx, errors.NewInternalServerError("保存文件失败", err))
		return
	}

	// 返回文件路径
	response.SuccessWithMessage(ctx, "上传成功", gin.H{
		"url": filePath,
	})
} 