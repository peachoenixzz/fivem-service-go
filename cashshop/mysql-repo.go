package cashshop

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

func (h Handler) getInitCashShop(c echo.Context, discordID string) (ResponseInitCashShop, error) {
	logger := mlog.Logg
	query := `
		SELECT count(u.id)
		     ,u.identifier
		     ,u.firstname  
		     ,u.lastname 
		     ,cp.point
		     ,v.expire_date
        FROM users as u 
        INNER JOIN cash_point as cp ON u.identifier = cp.discord_id
        INNER JOIN vip as v ON u.identifier = v.discord_id
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
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.Identifier, &res.FirstName, &res.LastName, &res.Point, &res.ExpireDateVip)
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

func (h Handler) GetCashShopItem(ctx context.Context, discordID string) ([]ResponseItemCashShop, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	var items []ResponseItemCashShop
	stmtStr := `
				SELECT
					DISTINCT ci.name AS item_name,
					ci.point AS price,
					ci.limit,
					ci.limit_type,
					CASE
						WHEN ci.limit_type = '00' THEN ci.limit
						WHEN ci.limit_type = '01' THEN ci.limit - COALESCE(daily_count.count, 0)
						WHEN ci.limit_type = '02' THEN ci.limit - COALESCE(hourly_count.count, 0)
						ELSE -1
					END AS remaining_quantity,
				    i.label as label_name
				FROM
				items i 
				INNER JOIN cash_items ci ON i.name = ci.name
				LEFT JOIN cash_history ch ON ci.name = ch.item_name AND ch.discord_id = ? AND DATE(ch.created_date) = CURDATE() --1
				LEFT JOIN (
					SELECT
						item_name,
						COUNT(1) AS count
					FROM
						cash_history
					WHERE
						discord_id = ? AND DATE(created_date) = CURDATE() --2
					GROUP BY
						item_name
				) AS daily_count ON ci.name = daily_count.item_name AND ci.limit_type = '01'
				LEFT JOIN (
					SELECT
						item_name,
						CASE
							WHEN (HOUR(created_date) BETWEEN 0 AND 5) THEN '0.00 - 6.00'
							WHEN (HOUR(created_date) BETWEEN 6 AND 11) THEN '6.00 - 12.00'
							WHEN (HOUR(created_date) BETWEEN 12 AND 17) THEN '12.00 - 18.00'
							WHEN (HOUR(created_date) BETWEEN 18 AND 23) THEN '18.00 - 0.00'
						END AS time_range,
						COUNT(1) AS count
					FROM
						cash_history
					WHERE
						discord_id = ? AND DATE(created_date) = CURDATE() --3
					GROUP BY
						item_name, time_range
				) AS hourly_count ON ci.name = hourly_count.item_name AND ci.limit_type = '02' 
					AND hourly_count.time_range = CASE
													  WHEN HOUR(NOW()) BETWEEN 0 AND 5 THEN '0.00 - 6.00'
													  WHEN HOUR(NOW()) BETWEEN 6 AND 11 THEN '6.00 - 12.00'
													  WHEN HOUR(NOW()) BETWEEN 12 AND 17 THEN '12.00 - 18.00'
													  WHEN HOUR(NOW()) BETWEEN 18 AND 23 THEN '18.00 - 0.00'
												  END
				ORDER BY
					ci.name`

	args := []interface{}{
		discordID,
		discordID,
		discordID,
	}

	stmt, err := h.MysqlDB.PrepareContext(ctx, stmtStr)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemCashShop{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemCashShop{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	for rows.Next() {
		var item ResponseItemCashShop
		err := rows.Scan(&item.Name, &item.Point, &item.MaxLimit, &item.LimitType, &item.RemainQuantity, &item.LabelName)
		if err != nil {
			logger.Error("Database Error : ", zap.Error(err))
			return []ResponseItemCashShop{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return []ResponseItemCashShop{}, echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	return items, nil
}

func (h Handler) ValidatePurchaseItem(tx *sql.Tx, req RequestBuyItem, discordID string) (ResponseValidateItem, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `
			SELECT
				DISTINCT ci.name AS item_name,
				ci.point AS price,
				ci.limit,
				ci.limit_type,
				CASE
					WHEN ci.limit_type = '00' THEN ci.limit
					WHEN ci.limit_type = '01' THEN ci.limit - COALESCE(daily_count.count, 0)
					WHEN ci.limit_type = '02' THEN ci.limit - COALESCE(hourly_count.count, 0)
					ELSE -1
				END AS remaining_quantity
			FROM
				cash_items ci
			LEFT JOIN cash_history ch ON ci.name = ch.item_name AND ch.discord_id = ? AND DATE(ch.created_date) = CURDATE()
			LEFT JOIN (
				SELECT
					item_name,
					COUNT(1) AS count
				FROM
					cash_history
				WHERE
					discord_id = ? AND DATE(created_date) = CURDATE()
				GROUP BY
					item_name
			) AS daily_count ON ci.name = daily_count.item_name AND ci.limit_type = '01'
			LEFT JOIN (
				SELECT
					item_name,
					CASE
						WHEN (HOUR(created_date) BETWEEN 0 AND 5) THEN '0.00 - 6.00'
						WHEN (HOUR(created_date) BETWEEN 6 AND 11) THEN '6.00 - 12.00'
						WHEN (HOUR(created_date) BETWEEN 12 AND 17) THEN '12.00 - 18.00'
						WHEN (HOUR(created_date) BETWEEN 18 AND 23) THEN '18.00 - 0.00'
					END AS time_range,
					COUNT(1) AS count
				FROM
					cash_history
				WHERE
					discord_id = ? AND DATE(created_date) = CURDATE()
				GROUP BY
					item_name, time_range
			) AS hourly_count ON ci.name = hourly_count.item_name AND ci.limit_type = '02' 
				AND hourly_count.time_range = CASE
												  WHEN HOUR(NOW()) BETWEEN 0 AND 5 THEN '0.00 - 6.00'
												  WHEN HOUR(NOW()) BETWEEN 6 AND 11 THEN '6.00 - 12.00'
												  WHEN HOUR(NOW()) BETWEEN 12 AND 17 THEN '12.00 - 18.00'
												  WHEN HOUR(NOW()) BETWEEN 18 AND 23 THEN '18.00 - 0.00'
											  END
			WHERE ci.name = ?
			ORDER BY
				ci.name;
			`

	args := []interface{}{
		discordID,
		discordID,
		discordID,
		req.Name,
	}

	row := tx.QueryRow(stmtStr, args...)

	var res ResponseValidateItem
	err := row.Scan(&res.Name, &res.Point, &res.MaxLimit, &res.LimitType, &res.RemainQuantity)
	if err != nil {
		tx.Rollback()
		logger.Error("Database Error : ", zap.Error(err))
		return ResponseValidateItem{}, err
	}

	return res, nil
}

func (h Handler) PurchaseItem(tx *sql.Tx, req RequestBuyItem, discordID string) (int64, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `
			UPDATE cash_point AS cp
			JOIN cash_items AS ci
			ON cp.discord_id = ? AND ci.name = ?
			SET cp.point = cp.point - ci.point
			WHERE cp.point >= ci.point;
	`

	args := []interface{}{
		discordID,
		req.Name,
	}

	r, err := tx.Exec(stmtStr, args...)
	if err != nil {
		// If there is an error, rollback the transaction
		tx.Rollback()
		logger.Error("Failed to Update record:", zap.Error(err))
		return 0, err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to retrieve affected rows: ", zap.Error(err))
		return 0, err
	}
	logger.Info("Row Affected Update Cash Point table", zap.Int64("row affected", rowsAffected))

	return rowsAffected, nil
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
