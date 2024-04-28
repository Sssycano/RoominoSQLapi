package api

import (
	"net/http"
	"roomino/service"
	"roomino/types"

	"github.com/gin-gonic/gin"
)

func UnitInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UnitInforeq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.GetAvailableUnitsWithPetPolicy(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdatePetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdatePets
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.UpdatePetInfo(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func CreatePetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetPets
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.CreatePet(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetPetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.GetPet(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
func GetInterestsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UnitRentIDReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.GetInterests(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
func CreateInterestsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.PostInterestReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.CreateInterests(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetComplexUnitinfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UnitRentIDReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.GetComplexUnitinfo(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func SearchInterestswithcondHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.InteresCondReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		taskSrv := service.GetTaskSrv()
		resp, err := taskSrv.SearchInterestswithcond(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func AveragePriceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AveragePriceReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		taskSrv := service.TaskSrv{}
		resp, err := taskSrv.GetAveragePrice(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
