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

		spec := getPaginationSpec(ctx.Request)

		res, err := vendorSvc.GetAll(ctx, spec)
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

func getPaginationSpec(r *http.Request) database.PaginationSpec {
	queryParam := r.URL.Query()

	pageString := queryParam.Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

	order := queryParam.Get("order")

	spec := database.PaginationSpec{
		Limit: 10,
		Order: order,
		Page:  page,
	}

	return spec
}
