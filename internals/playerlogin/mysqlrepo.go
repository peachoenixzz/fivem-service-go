package playerlogin

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) PlayerIdentify(c echo.Context, req Request) (Response, error) {
	logger := mlog.L(c)
	logger.Info("prepare to make query PlayerIdentify")
	query := "SELECT `identifier`,`job`,`group` FROM users WHERE identifier = ?"

	// Create a prepared statement
	logger.Info("mysql prepare query PlayerIdentify", zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	stmt, err := h.MysqlDB.Prepare(query)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err), zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	var res Response
	logger.Info("query Row PlayerIdentify", zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	err = stmt.QueryRow(req.Identifier).Scan(&res.Identifier, &res.Job, &res.Group)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err), zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
		return res, err
	}
	logger.Info("after query row and ready to return PlayerIdentify", zap.String("service", "playerlogin"), zap.String("discordID", req.Identifier))
	return res, nil
}
