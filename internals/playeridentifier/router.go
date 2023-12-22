package playeridentifier

import (
	"database/sql"
	"github.com/kkgo-software-engineering/workshop/config"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, mongodb *mongo.Client, mysqlDB *sql.DB) *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BasicAuth(mw.Authenicate()))
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)
	// Login route
	e.GET("/discord/id/:discordID", h.GetPlayerDiscordID)
	e.PUT("/discord/id/:discordID", h.UpdateVIPPoint)
	e.PUT("/cash", h.UpdateCashPointEndPoint)

	return e
}
