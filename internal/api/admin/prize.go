package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) prizeRoutes() {
	prizesGroup := s.engine.Group("/prizes")

	// swagger:route GET /prizes admin prizes getPrizes
	//
	// responses:
	//   200: GetPrizesResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	prizesGroup.GET("", func(c *gin.Context) {
		prizes, err := s.prizeService.GetPrizes(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, prizes)
	})
}
