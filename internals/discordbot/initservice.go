package discordbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	err = mysqlDB.Ping()
	if err != nil {
		fmt.Printf("mysql unable to ping, error: %v", err)
		os.Exit(1)
	}

	h := New(cfg.FeatureFlag, mysqlDB)

	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a WebSocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!ei" {
			h.handleEICommand(s, m)
		} else if m.Content == "!ev" {
			h.handleEVCommand(s, m)
		}
	})

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()

}
