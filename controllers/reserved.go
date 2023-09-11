package controllers

import (
	"fmt"

	argusapi "github.com/boardware-cloud/argus-api"
	services "github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
)

type ReservedApi struct{}

// ListReserved implements argusapi.ReservedApiInterface.
func (ReservedApi) ListReserved(gin_context *gin.Context) {
	panic("unimplemented")
}

// CreateReservedMonitor implements argusapi.ReservedApiInterface.
func (ReservedApi) CreateReservedMonitor(ctx *gin.Context, request argusapi.CreateReservedRequest) {
	middleware.IsRoot(ctx,
		func(c *gin.Context, account model.Account) {
			fmt.Println(request)
			services.CreateReserved(
				utils.StringToUint(*request.AccountId),
				*request.StartAt,
				*request.ExpiredAt,
			)
		})
}
