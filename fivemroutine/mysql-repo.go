package fivemroutine

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

func (h Handler) UpdateExpireVip(tx *sql.Tx) error {
	logger := mlog.Logg
	logger.Info("prepare to make Update expire vip by routine Go")
	stmtStr := `
			UPDATE vip v
			INNER JOIN
			(
				SELECT discord_id
				FROM vip
				WHERE SYSDATE() > DATE(expire_date) 
				AND priority != 'Citizen'
			) as select_expire_id
			ON
				v.discord_id  =  select_expire_id.discord_id
			SET
				v.vip_point  = 0 ,
				v.extra_point = 0 ,
			    v.priority = 'Citizen',
				v.total_point = permanant_point,
			    v.last_updated = SYSDATE()
`
	r, err := tx.Exec(stmtStr)
	if err != nil {
		// If there is an error, rollback the transaction
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to retrieve affected rows: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("Row Affected Update vip table", zap.Int64("row affected", rowsAffected))
	return nil
}
