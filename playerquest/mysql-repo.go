package playerquest

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

func (h Handler) QueryRequireQuest(c echo.Context, discordID string) (ResponseRequireQuestPlayer, error) {
	logger := mlog.Logg
	query := `
		SELECT count(tpw.discord_id),tlw.require , tpw.weight_level
              FROM TB_PLAYER_WEIGHT tpw
              INNER JOIN TB_LEVEL_WEIGHT tlw
 			  WHERE tlw.weight_level = tpw.weight_level +1 
              AND discord_id = ?
`

	// Create a prepared statement
	logger.Info("mysql prepare query Discord ID")
	stmt, err := h.MysqlDB.Prepare(query)
	if err != nil {
		logger.Error("sql error", zap.Error(err))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	logger.Info("query Row Discord Id")
	var rowCount int
	var res ResponseRequireQuestPlayer
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.Require, &res.WeightLevel)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return ResponseRequireQuestPlayer{}, err
	}
	logger.Info("after query row and ready to return Data")
	if rowCount == 0 {
		return ResponseRequireQuestPlayer{}, c.JSON(http.StatusNotFound, Message{Message: "not found discord id"})
	}
	return res, nil
}
