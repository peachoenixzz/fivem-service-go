package middleware

import (
	"crypto/subtle"
	"database/sql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Identifier string `json:"identifier"`
	Job        string `json:"job"`
	Group      string `json:"group"`
	jwt.RegisteredClaims
}

type Response struct {
	Identifier string `json:"identifier"`
	Job        string `json:"job"`
	Group      string `json:"group"`
}

type Handler struct {
	Cfg        config.FeatureFlag
	PostgresDB *sql.DB
	MongoDB    *mongo.Client
	mysqlDB    *sql.DB
}

type LoginService interface {
	Login()
}

func Authenicate() func(username, password string, c echo.Context) (bool, error) {
	return func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("thecircledev")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("Thecircle112!@#")) == 1 {
			return true, nil
		}
		return false, nil
	}
}

func LoginSuccess(c echo.Context, res Response) error {
	// Set custom claims
	logger := mlog.L(c)
	logger.Info("prepare JWT Data")
	claims := &JwtCustomClaims{
		res.Identifier,
		res.Job,
		res.Group,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	logger.Info("prepare encode HS256")
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	logger.Info("generate token for user")
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("550076b5-532c-439e-92d9-655f8207fdee"))
	if err != nil {
		logger.Error("failed to generate", zap.Error(err))
		return err
	}
	logger.Info("ready to return token to user")
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
