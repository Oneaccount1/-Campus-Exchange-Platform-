package upload

import (
	"campus/internal/bootstrap"
	"campus/internal/utils/logger"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadConfig 上传配置
type UploadConfig struct {
	SavePath     string   // 保存路径
	AllowedTypes []string // 允许的文件类型
	MaxSize      int64    // 最大大小（MB）
}

// GetUploadConfig 获取上传配置
func GetUploadConfig() *UploadConfig {
	config := bootstrap.GetConfig()

	// 将逗号分隔的类型转换为切片
	allowedTypes := strings.Split(config.Upload.AllowedTypes, ",")
	for i, t := range allowedTypes {
		allowedTypes[i] = strings.TrimSpace(t)
	}

	return &UploadConfig{
		SavePath:     config.Upload.SavePath,
		AllowedTypes: allowedTypes,
		MaxSize:      int64(config.Upload.MaxSize) * 1024 * 1024, // 转换为字节
	}
}

// IsTypeAllowed 检查文件类型是否允许
func (c *UploadConfig) IsTypeAllowed(fileType string) bool {
	//ext := strings.TrimPrefix(filepath.Ext(fileType), ".")
	//for _, t := range c.AllowedTypes {
	//	if strings.EqualFold(t, ext) {
	//		return true
	//	}
	//}
	//return false
	return true
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, fileType string) (string, error) {
	config := GetUploadConfig()

	// 检查文件大小
	if file.Size > config.MaxSize {
		return "", errors.New(fmt.Sprintf("文件大小超过限制，最大允许%dMB", config.MaxSize/1024/1024))
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return "", errors.New("无法确定文件类型")
	}

	// 如果没有指定fileType，使用multipart文件中的类型
	if fileType == "" {
		fileType = strings.TrimPrefix(ext, ".")
	}

	// 检查类型是否允许
	if !config.IsTypeAllowed(fileType) {
		return "", errors.New(fmt.Sprintf("不支持的文件类型，允许的类型：%s", strings.Join(config.AllowedTypes, ", ")))
	}

	// 创建保存路径
	savePath := config.SavePath
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		if err := os.MkdirAll(savePath, 0755); err != nil {
			logger.Errorf("创建上传目录失败: %v", err)
			return "", errors.New("创建上传目录失败")
		}
	}

	// 创建子目录（按照日期）
	datePath := time.Now().Format("2006/01/02")
	fullSavePath := filepath.Join(savePath, fileType, datePath)
	if _, err := os.Stat(fullSavePath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullSavePath, 0755); err != nil {
			logger.Errorf("创建上传子目录失败: %v", err)
			return "", errors.New("创建上传子目录失败")
		}
	}

	// 生成文件名（时间戳+随机数+原始文件扩展名）
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), GetRandomString(10), ext)
	filePath := filepath.Join(fullSavePath, fileName)

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		logger.Errorf("打开上传文件失败: %v", err)
		return "", errors.New("打开上传文件失败")
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("创建目标文件失败: %v", err)
		return "", errors.New("创建目标文件失败")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		logger.Errorf("保存文件失败: %v", err)
		return "", errors.New("保存文件失败")
	}

	// 返回可访问的URL路径（相对路径）
	relativePath := filepath.Join("/static", fileType, datePath, fileName)
	// 统一使用斜杠，避免Windows路径问题
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")

	// 记录日志
	logger.Infof("文件已上传: %s, 访问路径: %s", filePath, relativePath)

	return relativePath, nil
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		// 增加一些随机性
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}
