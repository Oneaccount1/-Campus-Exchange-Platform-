package controllers

import (
	"campus/internal/models"
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
	var newProduct models.Product
	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := c.service.CreateProduct(newProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, product)
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
	id := ctx.Param("id")
	var updatedProduct models.Product
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := c.service.UpdateProduct(id, updatedProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
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

// SearchProductsByKeyword 通过商品名称模糊查询商品
func (c *ProductController) SearchProductsByKeyword(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	products, err := c.service.SearchProductsByKeyword(keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}
