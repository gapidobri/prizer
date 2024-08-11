package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) wonPrizeRoutes() {
	wonPrizesGroup := s.engine.Group("/won-prizes")

	// swagger:route GET /won-prizes admin wonPrizes getWonPrizes
	//
	// parameters:
	//   + name: gameId
	//     in: query
	//     type: string
	//   + name: userId
	//     in: query
	//     type: string
	//   + name: prizeId
	//     in: query
	//     type: string
	//
	// responses:
	//   200: GetWonPrizesResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	wonPrizesGroup.GET("", func(c *gin.Context) {
		var filter api.GetWonPrizesFilter
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		wonPrizes, err := s.wonPrizeService.GetWonPrizes(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, wonPrizes)
	})
}
