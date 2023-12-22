package playerlogs

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

func RegRoute(cfg config.Config, logger *zap.Logger, postgresDB *sql.DB, mongodb *mongo.Client) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BasicAuth(mw.Authenicate()))
	hFiveMLog := New(cfg.FeatureFlag, postgresDB, mongodb)
	e.POST("/", hFiveMLog.AddFiveMLogEndPoint)
	e.POST("/custom", hFiveMLog.CustomLogEndPoint)
	e.GET("/", hFiveMLog.GetFiveMLogEndPoint)
	e.GET("/steamid/:steamid/events/:event", hFiveMLog.CaseEventAndSteamIDEndPoint)
	e.GET("/steamid/:steamid/events", hFiveMLog.AllEventAndSteamIDEndPoint)
	e.GET("/events/:event", hFiveMLog.ByEventEndPoint)

	return e
}
