package er

import (
	"github.com/gapidobri/prizer/pkg/errors"
	"net/http"
)

var (
	NotFound = errors.New(http.StatusNotFound, "not_found", "Not found")

	BadRequest = errors.New(http.StatusBadRequest, "invalid_body", "Invalid body")
)
