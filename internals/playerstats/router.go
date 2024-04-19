package playerstats

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
	e.Use(mw.RequestMetadataMiddleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BasicAuth(mw.Authenicate()))
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)

	// Login route
	e.GET("/money", h.GetAllMoney)
	e.GET("/playeritems", h.GetItemPlayer)
	e.GET("/playervault", h.GetItemVault)
	e.GET("/vehicle", h.GetVehicleByModel)
	//r.GET("auth", mw.Restricted)
	//hFiveMLog := New(cfg.FeatureFlag, postgresDB, mongodb)

	//r.POST("/", hFiveMLog.AddPoliceLogEndPoint)
	//e.GET("/", hFiveMLog.GetFiveMLogEndPoint)
	//e.GET("/steamid/:steamid/events/:event", hFiveMLog.CaseEventAndSteamIDEndPoint)
	//e.GET("/policelogs/steamid/:steamid/events", hFiveMLog.AllEventAndSteamIDEndPoint)
	//e.GET("/policelogs/events/:event", hFiveMLog.ByEventEndPoint)

	return e
}
