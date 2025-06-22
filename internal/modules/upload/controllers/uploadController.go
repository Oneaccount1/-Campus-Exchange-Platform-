package controllers

import (
	"campus/internal/utils/errors"
	"campus/internal/utils/logger"
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
	logger.Info("开始处理图片上传请求")

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Errorf("获取上传文件失败: %v", err)
		response.HandleError(ctx, errors.NewBadRequestError("获取上传文件失败", err))
		return
	}

	logger.Infof("接收到文件上传: %s, 大小: %d 字节", file.Filename, file.Size)

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		logger.Errorf("不支持的文件类型: %s", ext)
		response.HandleError(ctx, errors.NewBadRequestError("只支持jpg/jpeg/png/gif格式的图片", nil))
		return
	}

	// 保存文件
	fileType := "images" // 文件类型子目录
	filePath, err := upload.SaveUploadedFile(file, fileType)
	if err != nil {
		logger.Errorf("保存文件失败: %v", err)
		response.HandleError(ctx, errors.NewInternalServerError("保存文件失败", err))
		return
	}

	logger.Infof("文件上传成功: %s -> %s", file.Filename, filePath)

	// 返回文件路径
	response.SuccessWithMessage(ctx, "上传成功", gin.H{
		"url":       filePath,
		"image_url": filePath, // 添加image_url字段，与ProductImage模型一致
	})
}

// UploadFile 上传通用文件
func (c *UploadController) UploadFile(ctx *gin.Context) {
	logger.Info("开始处理通用文件上传请求")

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Errorf("获取上传文件失败: %v", err)
		response.HandleError(ctx, errors.NewBadRequestError("获取上传文件失败", err))
		return
	}

	logger.Infof("接收到文件上传: %s, 大小: %d 字节", file.Filename, file.Size)

	// 检查文件类型参数
	fileType := ctx.DefaultQuery("type", "files")
	logger.Infof("文件类型: %s", fileType)

	// 保存文件
	filePath, err := upload.SaveUploadedFile(file, fileType)
	if err != nil {
		logger.Errorf("保存文件失败: %v", err)
		response.HandleError(ctx, errors.NewInternalServerError("保存文件失败", err))
		return
	}

	logger.Infof("文件上传成功: %s -> %s", file.Filename, filePath)

	// 返回文件路径
	response.SuccessWithMessage(ctx, "上传成功", gin.H{
		"url": filePath,
	})
}
