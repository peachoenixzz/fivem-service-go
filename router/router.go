package router

import (
	"database/sql"
	fivemlogs "github.com/kkgo-software-engineering/workshop/fivemlogs"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, postgresDB *sql.DB, mongodb *mongo.Client) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(postgresDB)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hFiveMLog := fivemlogs.New(cfg.FeatureFlag, postgresDB, mongodb)
	e.POST("/fivemlogs", hFiveMLog.AddFiveMLogEndPoint)
	e.GET("/fivemlogs", hFiveMLog.GetFiveMLogEndPoint)
	e.GET("/fivemlogs/steamid/:steamid/events/:event", hFiveMLog.CaseEventAndSteamIDEndPoint)
	e.GET("/fivemlogs/steamid/:steamid/events", hFiveMLog.AllEventAndSteamIDEndPoint)
	e.GET("/fivemlogs/events/:event", hFiveMLog.ByEventEndPoint)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	return e
}
