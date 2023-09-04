package cashshop

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
	"net/http"
	"sync"
)

var (
	clientLock  = &sync.Mutex{}
	clientFlags = make(map[string]bool)
)

func PerClientRateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()

		clientLock.Lock()
		if _, exists := clientFlags[ip]; exists {
			clientLock.Unlock()
			return c.JSON(http.StatusTooManyRequests, map[string]string{"message": "already processing"})
		}
		// Set the flag to indicate processing for this client
		clientFlags[ip] = true
		clientLock.Unlock()

		// Make sure to unset the flag after the request is done or times out
		defer func() {
			clientLock.Lock()
			delete(clientFlags, ip)
			clientLock.Unlock()
		}()

		// Call the next handler in the chain
		return next(c)
	}
}

func RegRoute(cfg config.Config, logger *zap.Logger, mongodb *mongo.Client, mysqlDB *sql.DB) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)

	JWTConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(mw.JwtCustomClaims)
		},
		SigningKey: []byte("550076b5-532c-439e-92d9-655f8207fdee"),
	}
	e.Use(echojwt.WithConfig(JWTConfig))
	///cash-shop/
	e.GET("/users", h.GetInitCashShopEndPoint)
	e.GET("/users/items", h.GetCashShopItemEndPoint)
	e.PUT("/users/buy", h.BuyCashShopEndPoint, PerClientRateLimiter)
	return e
}
