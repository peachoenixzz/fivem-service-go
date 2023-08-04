package playeridentifier

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) InsertPlayerItems(tx *sql.Tx, i Item, id string) error {
	logger := mlog.Logg
	logger.Info("prepare to make Insert Vip give player items ")
	stmtStr := `
		INSERT
		INTO TB_GIVE_PLAYERS_ITEMS (item_name,quantity,identifier)
		VALUES (?,?,?)
	`
	args := []interface{}{
		i.ItemName,
		i.Quantity,
		id,
	}
	r, err := tx.Exec(stmtStr, args...)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Insert Vip give player items:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("done to Insert Vip give player items ")
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Insert Vip give player items: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("Row Affected Vip give player items", zap.Int64("row affected", rowsAffected))

	return nil
}

func (h Handler) InsertExpireDateItem(tx *sql.Tx, i Item, id string) error {
	logger := mlog.Logg
	logger.Info("prepare to do stmt expire date ")
	stmtStr := `
		INSERT
		INTO items_expire (item_name,player_id,expire_timestamp)
		VALUES (?,?,?)
	`
	logger.Info("prepare to calculate expire date ")
	currentTime := time.Now()
	expireDate := currentTime.AddDate(0, 0, i.ExpireDate)
	logger.Info("done to calculate expire date ")
	args := []interface{}{
		i.ItemName,
		id,
		expireDate,
	}
	logger.Info("prepare to Insert expire date ")
	r, err := tx.Exec(stmtStr, args...)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Insert Expire Item record:", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("done to Insert expire date ")
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to Insert Expire Item record: ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}
	logger.Info("Row Affected Update vip table", zap.Int64("row affected", rowsAffected))

	return nil
}

func (h Handler) UpdateVipPointByPlayerDiscord(tx *sql.Tx, req RequestUpdateVip, discordID string) error {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	stmtStr := `
	UPDATE vip 
		SET 
			vip_point = ?,
			extra_point = ?,
			permanant_point = ?,
			priority = ?,
			total_point = ?,
		    expire_date = 
		    CASE
				WHEN expire_date > DATE_ADD(SYSDATE(), INTERVAL 7 DAY) THEN expire_date
			ELSE DATE_ADD(SYSDATE(), INTERVAL 30 DAY)
		END,
			last_updated = SYSDATE()
		WHERE
			discord_id = ?
`
	//currentTime := time.Now()
	totalPoint := HandleTotalVipPoint(req.VipPoint, req.ExtraPoint, req.PermanentPoint)
	//expireDate := currentTime.AddDate(0, 0, 30)
	fmt.Println(req.VipPoint, req.ExtraPoint, req.PermanentPoint, req.Priority, discordID)
	args := []interface{}{
		req.VipPoint,
		req.ExtraPoint,
		req.PermanentPoint,
		req.Priority,
		totalPoint,
		//expireDate,
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

func (h Handler) QueryPlayerDiscord(c echo.Context, discordID string) (Response, error) {
	logger := mlog.Logg
	logger.Info("prepare to make query Discord ID")
	query := `
		SELECT 
			count(discord_id) AS countID,
			discord_id, 
			steam_id, 
			vip_point, 
			extra_point, 
			permanant_point, 
			priority,
			identifier,
			DATE_FORMAT(expire_date,'%Y-%m-%d') AS expire_date , 
			DATE_FORMAT(last_updated ,'%Y-%m-%d') AS last_updated
		FROM vip 
		WHERE discord_id = ?
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
	var res Response
	err = stmt.QueryRow(discordID).Scan(&rowCount, &res.DiscordID, &res.SteamID,
		&res.VipPoint, &res.ExtraPoint, &res.PermanentPoint, &res.Priority, &res.Identifier, &res.ExpireDate,
		&res.LastUpdated)
	if err != nil {
		logger.Error("query row fail ", zap.Error(err))
		return Response{}, err
	}
	logger.Info("after query row and ready to return PlayerIdentify")
	if rowCount == 0 {
		return Response{}, c.JSON(http.StatusNotFound, Message{Message: "not found discord id"})
	}
	return res, nil
}
