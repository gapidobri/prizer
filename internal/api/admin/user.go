package admin

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) userRoutes() {
	usersGroup := s.engine.Group("/users")

	// swagger:route GET /users admin users getUsers
	//
	// parameters:
	//   + name: gameId
	//     in: query
	//     type: string
	//
	// responses:
	//   200: GetUsersResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	usersGroup.GET("", func(c *gin.Context) {
		var filter api.GetUsersFilter
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		users, err := s.userService.GetUsers(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, users)
	})

	// swagger:route GET /users/{userId} admin users getUser
	//
	// parameters:
	//   + name: userId
	//     in: path
	//     type: string
	//     required: true
	//
	// responses:
	//   200: GetUserResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	usersGroup.GET(":userId", func(c *gin.Context) {
		response, err := s.userService.GetUser(c.Request.Context(), c.Param("userId"))
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
