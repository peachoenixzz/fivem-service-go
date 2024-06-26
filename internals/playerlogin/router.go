package playerlogin

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kkgo-software-engineering/workshop/config"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	echojwt "github.com/labstack/echo-jwt/v4"
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
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)
	// Login route
	e.POST("/users", h.GetPlayerIdentify)
	JWTConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(mw.JwtCustomClaims)
		},
		SigningKey: []byte("550076b5-532c-439e-92d9-655f8207fdee"),
	}
	r := e.Group("/")
	r.Use(echojwt.WithConfig(JWTConfig))

	return e
}
