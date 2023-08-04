package playeritems

import (
	"context"
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

func (h Handler) AllVipItems(ctx context.Context) ([]ResponseVipItems, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query AllItems")
	stmtStr := `
		SELECT 
			name,
			label
		FROM 
			items
		WHERE 
			type = "vip"
	`

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseVipItems{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseVipItems{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	var items []ResponseVipItems
	for rows.Next() {
		var res ResponseVipItems
		err := rows.Scan(&res.ItemName, &res.ItemLabel)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponseVipItems{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, res)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseVipItems{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	return items, nil
}
