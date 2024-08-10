package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) userRoutes() {
	usersGroup := s.engine.Group("/users")

	// swagger:route GET /users admin users getUsers
	//
	// responses:
	//   200: GetUsersResponse
	//   400: ErrorResponse
	//   403: ErrorResponse
	//   500: ErrorResponse
	//
	usersGroup.GET("", func(c *gin.Context) {
		users, err := s.userService.GetUsers(c.Request.Context())
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
