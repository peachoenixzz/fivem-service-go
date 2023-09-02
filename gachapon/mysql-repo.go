package gachapon

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}
func (h Handler) getInitGachapon(c echo.Context, discordID string) (ResponseInitGacha, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	query := `SELECT count(firstname),firstname , lastname
             FROM users 
             WHERE users.identifier = ?;
`

	// Create a prepared statement
	logger.Info("mysql prepare query Discord ID")
	stmt, err := h.MysqlDB.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	logger.Info("query Row Discord Id")
	var rowCount int
	var res ResponseInitGacha
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.Name, &res.LabelName)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return ResponseInitGacha{}, err
	}
	logger.Info("after query row and ready to return InitGachapon")
	if rowCount == 0 {
		return ResponseInitGacha{}, c.JSON(http.StatusNotFound, Message{Message: "not found discord id"})
	}
	return res, nil
}
