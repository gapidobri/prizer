package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) gameRoutes() {
	game := s.engine.Group("/games")

	// swagger:route GET /games admin games getGames
	//
	// responses:
	//   200: GetGamesResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	game.GET("", func(c *gin.Context) {
		response, err := s.gameService.GetGames(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response)
	})

	// swagger:route GET /games/{gameId} admin games getGame
	//
	// responses:
	//   200: GetGameResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	game.GET("", func(c *gin.Context) {
		response, err := s.gameService.GetGames(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
