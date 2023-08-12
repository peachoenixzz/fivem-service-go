package cashshop

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

func (h Handler) getInitCashShop(c echo.Context, discordID string) (ResponseInitCashShop, error) {
	logger := mlog.Logg
	query := `
		SELECT count(u.id),u.identifier, u.firstname  , u.lastname , cp.point
        FROM users as u 
        INNER JOIN cash_point as cp ON u.identifier = cp.discord_id
		WHERE u.identifier = ?
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
	var res ResponseInitCashShop
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.Identifier, &res.FirstName, &res.LastName, &res.Point)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return ResponseInitCashShop{}, err
	}
	logger.Info("after query row and ready to return PlayerIdentify")
	if rowCount == 0 {
		return ResponseInitCashShop{}, c.JSON(http.StatusNotFound, Message{Message: "not found discord id"})
	}
	return res, nil
}
func (h Handler) UpdateCashPoint(tx *sql.Tx, req RequestUpdatePoint, discordID string) error {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `
	UPDATE cash_point 
		SET 
			point = point + ?
		WHERE
			discord_id = ?
`

	args := []interface{}{
		req.CashPoint,
		discordID,
	}
	r, err := tx.Exec(stmtStr, args...)
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
