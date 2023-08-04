package cashshop

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
)

type Message struct {
	Status  int
	Message interface{}
}

type ResponseInitCashShop struct {
}

type Request struct {
}

func (h Handler) GetInitCashShopEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	fmt.Println("JOB", playerInfo.Job)
	fmt.Println("identifier", playerInfo.Identifier)
	fmt.Println("group", playerInfo.Group)
	res, err := h.getInitCashShop(context.Background(), playerInfo)
	if err != nil {
		logger.Error("got error when query DB : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "query error")
	}
	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}
