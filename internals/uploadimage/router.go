package uploadimage

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kkgo-software-engineering/workshop/config"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, mongodb *mongo.Client, mysqlDB *sql.DB) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)
	r := e.Group("/")

	JWTConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(mw.JwtCustomClaims)
		},
		SigningKey: []byte("550076b5-532c-439e-92d9-655f8207fdee"),
	}

	r.Use(echojwt.WithConfig(JWTConfig))
	r.POST("create", h.AddImageEndPoint)

	return e
}
