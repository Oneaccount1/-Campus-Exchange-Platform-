package controllers

import (
	"campus/internal/modules/user/api"
	"campus/internal/modules/user/services"
	"campus/internal/utils/errors"
	"campus/internal/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FavoriteController struct {
	service services.FavoriteService
}

func (f *FavoriteController) AddFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}
	var req api.FavoriteRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	if err := f.service.AddFavorite(userID.(uint), req.ProductID); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "收藏成功", nil)
}

func (f *FavoriteController) RemoveFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效的商品ID", err))
		return
	}
	if err = f.service.RemoveFavorite(userID.(uint), uint(productID)); err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.SuccessWithMessage(ctx, "收藏成功", nil)
}

func (f *FavoriteController) ListFavorites(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req api.QueryPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}
	favorites, err := f.service.ListUserFavorites(userID.(uint), req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}
	response.Success(ctx, favorites)
}

// CheckFavorite 检查是否收藏
func (f *FavoriteController) CheckFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		response.HandleError(ctx, errors.NewBadRequestError("无效的商品ID", err))
		return
	}
	isFavorite, err := f.service.CheckIsFavorite(userID.(uint), uint(productID))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{"is_favorite": isFavorite})

}

func (f *FavoriteController) ListUserProducts(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req api.QueryPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.HandleError(ctx, errors.NewValidationError("请求参数错误", err))
		return
	}

	products, err := f.service.GetUserProducts(userID.(uint), req.Page, req.Size)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, products)
}

// NewFavoriteController 创建收藏控制器
func NewFavoriteController() *FavoriteController {
	return &FavoriteController{
		service: services.NewFavoriteService(),
	}
}
