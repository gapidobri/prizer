package api

import (
	"github.com/elliotxx/healthcheck"
	"github.com/elliotxx/healthcheck/checks"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewServer(db *sqlx.DB) *gin.Engine {
	engine := gin.Default()

	_ = healthcheck.Register(&engine.RouterGroup)

	engine.GET("/readyz", healthcheck.NewHandler(
		healthcheck.NewDefaultHandlerConfigFor(checks.NewSQLCheck(db.DB))),
	)

	engine.Use(ErrorHandler)

	return engine
}
