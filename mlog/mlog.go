package mlog

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
)

const key = "logger"

var Logg *zap.Logger

func init() {
	var err error
	Logg, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
}

func L(c echo.Context) *zap.Logger {
	switch logger := c.Get(key).(type) {
	case *zap.Logger:
		return logger
	default:
		return zap.NewNop()
	}
}
