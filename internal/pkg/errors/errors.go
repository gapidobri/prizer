package er

import (
	"github.com/gapidobri/prizer/pkg/errors"
	"net/http"
)

var (
	AlreadyRolled               = errors.New(http.StatusConflict, "already_rolled", "You have already rolled today")
	GameNotFound                = errors.New(http.StatusNotFound, "game_not_found", "Game not found")
	CollaboratorExists          = errors.New(http.StatusConflict, "collaborator_exists", "Collaborator already exists")
	CollaborationMethodNotFount = errors.New(http.StatusNotFound, "collaboration_method_not_found", "Collaboration method not found")
)
