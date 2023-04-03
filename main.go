package main

import (
	"github.com/kkgo-software-engineering/workshop/playerlogin"
	"github.com/kkgo-software-engineering/workshop/policelogs"
	"log"
	"os"

	"github.com/kkgo-software-engineering/workshop/playerlogs"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("playerlogin") != "" {
		logger.Info("prepare to playerlogin")
		playerlogin.InitService()
		logger.Info("Registered FiveM log service on /playerlogs")
	}

	if os.Getenv("policelogs") != "" {
		logger.Info("prepare to policelogs")
		policelogs.InitService()
		logger.Info("Registered FiveM log service on /policelogs")
	}

	//fmt.Println(os.Getenv("playerlogs"))
	// Register the FiveM log service if enabled
	if os.Getenv("playerlogs") != "" {
		logger.Info("prepare to playerlogs")
		playerlogs.InitService()
		logger.Info("Registered FiveM log service on /playerlogs")
	}

	logger.Info("Register service fail")
	//// Register the FiveM police log service if enabled
	//if *fivempolicelogEnabled {
	//	fivempolicelog := NewFiveMPoliceLogService()
	//	mux.HandleFunc("/fivempolicelog", fivempolicelog.HandleRequest)
	//	log.Printf("Registered FiveM police log service on /fivempolicelog")
	//}

}
