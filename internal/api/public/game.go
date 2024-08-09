package public

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) gameRoutes() {
	games := s.engine.Group("/games")

	// swagger:route POST /games/{gameId}/roll games roll
	//
	// parameters:
	//   + name: gameId
	//     in: path
	//     type: string
	//   + name: body
	//     in: body
	//     type: RollRequest
	//
	// responses:
	//   200: RollResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	games.POST(":gameId/roll", func(c *gin.Context) {
		var rollRequest api.RollRequest
		err := c.ShouldBind(&rollRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res, err := s.gameService.Roll(c.Request.Context(), c.Param("gameId"), rollRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}
