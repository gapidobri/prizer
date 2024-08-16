package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) prizeRoutes() {
	group := s.engine.Group("/prizes")

	// swagger:route GET /prizes admin prizes getPrizes
	//
	// parameters:
	//   + name: gameId
	//     in: query
	//     type: string
	//
	// responses:
	//   200: GetPrizesResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.GET("", func(c *gin.Context) {
		var filter api.GetPrizesFilter
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		prizes, err := s.prizeService.GetPrizes(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, prizes)
	})

	// swagger:route POST /prizes admin prizes createPrize
	//
	// parameters:
	//   + name: body
	//     in: body
	//     type: CreatePrizeRequest
	//
	// responses:
	//   201: EmptyResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.POST("", func(c *gin.Context) {
		var prize api.CreatePrizeRequest
		err := c.ShouldBindJSON(&prize)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = s.prizeService.CreatePrize(c.Request.Context(), prize)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, prize)
	})

	// swagger:route DELETE /prizes/{prizeId} admin prizes deletePrize
	//
	// parameters:
	//   + name: prizeId
	//     in: path
	//     type: string
	//     required: true
	//
	// responses:
	//   204: EmptyResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.DELETE(":prizeId", func(c *gin.Context) {
		prizeId := c.Param("prizeId")

		err := s.prizeService.DeletePrize(c.Request.Context(), prizeId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})
}
