package er

import (
	"github.com/gapidobri/prizer/pkg/errors"
	"net/http"
)

var (
	AlreadyParticipated         = errors.New(http.StatusConflict, "already_participated", "You have already participated today")
	GameNotFound                = NotFound.New("game_not_found", "Game not found")
	UserExists                  = errors.New(http.StatusConflict, "user_exists", "User already exists")
	UserNotFound                = NotFound.New("user_not_found", "User not found")
	ParticipationMethodNotFound = NotFound.New("participation_method_not_found", "Participation method not found")
	InvalidEmail                = errors.New(http.StatusBadRequest, "invalid_email", "Invalid email")
	ParticipationDataExists     = errors.New(http.StatusBadRequest, "participation_data_exists", "Participation data already exists")
)
