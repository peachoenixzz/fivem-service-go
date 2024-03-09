package uploadimage

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) UploadImage(c echo.Context, f *multipart.FileHeader) (Response, error) {
	logger := mlog.Logg
	logger.Info("prepare to open file image")

	src, err := f.Open()
	if err != nil {
		logger.Error("Error when Open file err (Open) : ", zap.Error(err))
		return Response{http.StatusInternalServerError, "", err.Error()}, err
	}

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			logger.Error("Error when Open file err (Close) : ", zap.Error(err))
		}
	}(src)

	u := c.Get("user").(*jwt.Token)
	p := u.Claims.(*mw.JwtCustomClaims)
	fmt.Println("Job : ", p.Job)
	fmt.Println("Identifier : ", p.Identifier)
	dir := fmt.Sprintf("images/%s-%d.jpg", p.Identifier, time.Now().Unix())
	dst, err := os.Create(dir)
	if err != nil {
		logger.Error("Failed to move file to destination err (Create) : ", zap.Error(err))
		return Response{http.StatusInternalServerError, "", err.Error()}, err
	}

	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			logger.Error("Failed to move file to destination err (Close) : ", zap.Error(err))
		}
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		logger.Error(fmt.Sprintf("Failed to move file to destination err (Create) : %s", err.Error()))
		return Response{http.StatusInternalServerError, "", err.Error()}, err
	}
	url := fmt.Sprintf("https://mongkol.dev/%s", dir)
	return Response{http.StatusOK, url, "Create Successfully"}, nil
}
