package uploadimage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitService() {

	cfg := config.New().All()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MySQL : ", cfg.MySQLDBConnection)
	mysqlDB, err := sql.Open("mysql", cfg.MySQLDBConnection)
	if err != nil {
		logger.Fatal("unable to configure database", zap.Error(err))
	}

	ctx := context.Background()
	fmt.Println("mongourl : ", cfg.MongoDBConnection)
	clientOptions := options.Client().ApplyURI(cfg.MongoDBConnection)
	mongoDB, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("unable to configure database", zap.Error(err))
	}

	e := RegRoute(cfg, logger, mongoDB, mysqlDB)
	err = mongoDB.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("mongodb unable to ping, error: %v", err)
		os.Exit(1)
	}

	err = mysqlDB.Ping()
	if err != nil {
		fmt.Printf("mysql unable to ping, error: %v", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)

	go func() {
		err := e.Start(addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("unexpected shutdown the server", zap.Error(err))
		}
		logger.Info("gracefully shutdown the server")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gCtx := context.Background()
	ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}
