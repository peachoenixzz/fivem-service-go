package fivemroutine

import (
	"database/sql"
	"fmt"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"time"
)

type VipDetail struct {
	DiscordID  string
	ExpireDate string
}

func initRoutine(cfg config.Config, logger *zap.Logger, mongodb *mongo.Client, mysqlDB *sql.DB) {
	h := New(cfg.FeatureFlag, mongodb, mysqlDB)
	logger.Info("prepare to init routine")
	if os.Getenv("routineTest") == "" {
		go h.Routine()
	}
	if os.Getenv("routineTest") != "" {
		go h.RoutineTest()
	}
	select {}
}

func (h Handler) Routine() {
	for {
		// Perform the recheck operation
		tx, err := h.MysqlDB.Begin()
		if err != nil {
			fmt.Println("Database Error : ", err)
			continue
		}
		err = h.UpdateExpireVip(tx)
		if err != nil {
			fmt.Println("Database Error : ", err)
			continue
		}
		err = tx.Commit()
		if err != nil {
			continue
		}
		// Sleep for 10 minutes
		time.Sleep(1 * time.Hour)
	}
}

func (h Handler) RoutineTest() {
	for {
		// Perform the recheck operation
		tx, err := h.MysqlDB.Begin()
		if err != nil {
			fmt.Println("Database Error : ", err)
			continue
		}
		err = h.UpdateExpireVip(tx)
		if err != nil {
			fmt.Println("Database Error : ", err)
			continue
		}
		err = tx.Commit()
		if err != nil {
			continue
		}
		// Sleep for 10 minutes
		time.Sleep(1 * time.Minute)
	}
}
