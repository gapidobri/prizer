package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) gameRoutes() {
	gamesGroup := s.engine.Group("/games")

	// swagger:route GET /games admin games getGames
	//
	// responses:
	//   200: GetGamesResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	gamesGroup.GET("", func(c *gin.Context) {
		response, err := s.gameService.GetGames(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response)
	})

	// swagger:route GET /games/{gameId} admin games getGame
	//
	// parameters:
	//   + name: gameId
	//     in: path
	//     type: string
	//
	// responses:
	//   200: GetGameResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	gamesGroup.GET(":gameId", func(c *gin.Context) {
		response, err := s.gameService.GetGame(c.Request.Context(), c.Param("gameId"))
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
