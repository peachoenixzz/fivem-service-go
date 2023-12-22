package uploadimage

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Response struct {
	Status   int    `json:"status"`
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

func (h Handler) AddImageEndPoint(c echo.Context) error {
	logger := mlog.L(c)

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*mw.JwtCustomClaims)

	fmt.Println("Job : ", claims.Job)
	fmt.Println("Identifier : ", claims.Identifier)
	fmt.Println("Group : ", claims.Group)
	fmt.Println("Claim : ", claims.RegisteredClaims)

	file, err := c.FormFile("file")
	//fmt.Println("form file", file)
	if err != nil {
		logger.Error("Failed upload image err : ", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	r, err := h.UploadImage(c, file)
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("Service Error : ", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create file successfully")
	return c.JSON(http.StatusOK, r)
}
