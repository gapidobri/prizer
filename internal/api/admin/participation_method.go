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
	//   + name: game_id
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

	// swagger:route PUT /participation-methods/{participationMethodId} admin participationMethods updateParticipationMethod
	//
	// parameters:
	//   + name: participationMethodId
	//     in: path
	//     type: string
	//     required: true
	//   + name: body
	//     in: body
	//     type: UpdateParticipationMethodRequest
	//     required: true
	//
	// responses:
	//   204: EmptyResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	group.PUT(":participationMethodId", func(c *gin.Context) {
		participationMethodId := c.Param("participationMethodId")

		var participationMethod api.UpdateParticipationMethodRequest
		err := c.ShouldBindJSON(&participationMethod)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = s.participationMethodService.UpdateParticipationMethod(
			c.Request.Context(), participationMethodId, participationMethod,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})

	drawMethods := group.Group(":participationMethodId/draw-methods")

	// swagger:route POST /participation-methods/{participationMethodId}/draw-methods/{drawMethodId} admin participationMethods linkDrawMethod
	//
	// parameters:
	//   + name: participationMethodId
	//     in: path
	//     type: string
	//     required: true
	//   + name: drawMethodId
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
	drawMethods.POST(":drawMethodId", func(c *gin.Context) {
		participationMethodId := c.Param("participationMethodId")
		drawMethodId := c.Param("drawMethodId")

		err := s.participationMethodService.LinkDrawMethod(c.Request.Context(), participationMethodId, drawMethodId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})

	// swagger:route DELETE /participation-methods/{participationMethodId}/draw-methods/{drawMethodId} admin participationMethods unlinkDrawMethod
	//
	// parameters:
	//   + name: participationMethodId
	//     in: path
	//     type: string
	//     required: true
	//   + name: drawMethodId
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
	drawMethods.DELETE(":drawMethodId", func(c *gin.Context) {
		participationMethodId := c.Param("participationMethodId")
		drawMethodId := c.Param("drawMethodId")

		err := s.participationMethodService.UnlinkDrawMethod(c.Request.Context(), participationMethodId, drawMethodId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})
}
