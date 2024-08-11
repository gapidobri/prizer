package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) prizeRoutes() {
	prizesGroup := s.engine.Group("/prizes")

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
	prizesGroup.GET("", func(c *gin.Context) {
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
}
