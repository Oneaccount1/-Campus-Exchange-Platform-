package controllers

import (
	"campus/internal/modules/product/api"
	"campus/internal/modules/product/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController() *ProductController {
	return &ProductController{
		service: services.NewProductService(),
	}
}

func (c *ProductController) ListProducts(ctx *gin.Context) {
	var req api.GetProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	products, err := c.service.GetAllProducts(req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, products)
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := c.service.GetProductByID(id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, product)
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req api.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	product, err := c.service.CreateProduct(&req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "商品创建成功", product)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var req api.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	product, err := c.service.UpdateProduct(id, &req)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "商品更新成功", product)
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteProduct(id); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "商品删除成功", nil)
}

func (c *ProductController) SearchProductsByKeyword(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	var req api.GetProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	products, err := c.service.SearchProductsByKeyword(keyword, req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, products)
}
