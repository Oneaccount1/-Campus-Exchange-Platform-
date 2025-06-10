package controllers

import (
	"campus/internal/modules/product/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ProductController 商品控制器结构体
type ProductController struct {
	service services.ProductService
}

// NewProductController 创建新的商品控制器实例
func NewProductController() *ProductController {
	return &ProductController{
		service: services.NewProductService(),
	}
}

// ListProducts 获取商品列表
func (c *ProductController) ListProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// CreateProduct 创建新商品
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	// 这里需要解析请求体
	// 示例代码，实际需要根据模型实现
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

// GetProductByID 根据ID获取商品
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := c.service.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// UpdateProduct 更新商品信息
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	//id := ctx.Param("id")
	// 这里需要解析请求体
	// 示例代码，实际需要根据模型实现
	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// DeleteProduct 删除商品
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteProduct(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
