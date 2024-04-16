package playerlogs

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitService() {

	cfg := config.New().All()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sqlDB : ", cfg.DBConnection)
	postgresDB, err := sql.Open("postgres", cfg.DBConnection)
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

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		logger.Fatal("error creating Discord session", zap.Error(err))
	}
	defer dg.Close()

	e := RegRoute(cfg, logger, postgresDB, mongoDB, dg)
	err = mongoDB.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("unable to ping, error: %v", err)
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
