package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) drawMethodRoutes() {
	group := s.engine.Group("/draw-methods")

	// swagger:route GET /draw-methods admin drawMethods getDrawMethods
	//
	// parameters:
	//   + name: gameId
	//     in: query
	//     type: string
	//   + name: participationMethodId
	//     in: query
	//     type: string
	//
	// responses:
	//   200: GetDrawMethodsResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.GET("", func(c *gin.Context) {
		var filter api.GetDrawMethodsFilter
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		drawMethods, err := s.drawMethodService.GetDrawMethods(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, drawMethods)
	})
}
