package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) participationMethodRoutes() {
	group := s.engine.Group("/participation-methods")

	// swagger:route GET /participation-methods admin participationMethods getParticipationMethods
	//
	// parameters:
	//   + name: gameId
	//     in: query
	//     type: string
	//
	// responses:
	//   200: GetParticipationMethodsResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.GET("", func(c *gin.Context) {
		var filter api.GetParticipationMethodsFilter
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		participationMethods, err := s.participationMethodService.GetParticipationMethods(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, participationMethods)
	})
}
