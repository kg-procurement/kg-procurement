package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/vendors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewVendorEngine(
	r *gin.Engine,
	cfg config.VendorRoutes,
	vendorSvc *vendors.VendorService,
) {
	r.GET(cfg.GetAll, func(ctx *gin.Context) {

		spec, err := getPaginationSpec(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		res, err := vendorSvc.GetAll(ctx, *spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetByLocation, func(ctx *gin.Context) {
		location := ctx.Query("location")

		res, err := vendorSvc.GetByLocation(ctx, location)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)
	})
}

func getPaginationSpec(r *http.Request) (*vendors.ServiceGetAllPaginationSpec, error) {
	queryParam := r.URL.Query()

	limitString := queryParam.Get("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return nil, err
	}

	pageString := queryParam.Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		return nil, err
	}

	order := database.ValidateOrderString(queryParam.Get("order"))

	spec := vendors.ServiceGetAllPaginationSpec{
		Limit: limit,
		Order: order,
		Page:  page,
	}

	return &spec, nil
}
