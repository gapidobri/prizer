package public

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) participationMethodRoutes() {
	method := s.engine.Group("/participationMethod")

	// swagger:route POST /participationMethod/{participationMethodId}/participate participationMethod participate
	//
	// parameters:
	//   + name: participationMethodId
	//     in: path
	//     type: string
	//   + name: body
	//     in: body
	//     type: ParticipationRequest
	//
	// responses:
	//   200: ParticipationResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	method.POST(":participationMethodId/participate", func(c *gin.Context) {
		var participationRequest api.ParticipationRequest
		err := c.ShouldBind(&participationRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		participationMethodId := c.Param("participationMethodId")

		res, err := s.gameService.Participate(c.Request.Context(), participationMethodId, participationRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}
