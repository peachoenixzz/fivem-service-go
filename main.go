package main

import (
	"github.com/kkgo-software-engineering/workshop/cashshop"
	"github.com/kkgo-software-engineering/workshop/fivemroutine"
	"github.com/kkgo-software-engineering/workshop/playeridentifier"
	"github.com/kkgo-software-engineering/workshop/playeritems"
	"github.com/kkgo-software-engineering/workshop/playerlogin"
	"github.com/kkgo-software-engineering/workshop/playerstats"
	"github.com/kkgo-software-engineering/workshop/uploadimage"
	"os"

	"github.com/kkgo-software-engineering/workshop/playerlogs"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()

	//os.Setenv("cashshop", "true")

	if os.Getenv("cashshop") != "" {
		logger.Info("prepare to cashshop")
		cashshop.InitService()
		logger.Info("Registered FiveM log service on /cash-shop")
	}

	if os.Getenv("playerlogin") != "" {
		logger.Info("prepare to playerlogin")
		playerlogin.InitService()
		logger.Info("Registered FiveM log service on /playerlogin")
	}

	//fmt.Println(os.Getenv("playerlogs"))
	// Register the FiveM log service if enabled
	if os.Getenv("playerlogs") != "" {
		logger.Info("prepare to playerlogs")
		playerlogs.InitService()
		logger.Info("Registered FiveM log service on /playerlogs")
	}

	if os.Getenv("uploadimages") != "" {
		logger.Info("prepare to uploadimages")
		uploadimage.InitService()
		logger.Info("Registered FiveM log service on /uploadimages")
	}

	if os.Getenv("vip") != "" {
		logger.Info("prepare to init vip")
		playeridentifier.InitService()
		logger.Info("Registered FiveM service on /vip")
	}

	if os.Getenv("items") != "" {
		logger.Info("prepare to init items")
		playeritems.InitService()
		logger.Info("Registered FiveM service on /items")
	}

	if os.Getenv("routine") != "" {
		logger.Info("prepare to init routine")
		fivemroutine.InitService()
		logger.Info("Registered FiveM service on /routine")
	}
	if os.Getenv("playerstats") != "" {
		logger.Info("prepare to init routine")
		playerstats.InitService()
		logger.Info("Registered FiveM service on /routine")
	}

	logger.Info("Register service fail")
	//// Register the FiveM police log service if enabled
	//if *fivempolicelogEnabled {
	//	fivempolicelog := NewFiveMPoliceLogService()
	//	mux.HandleFunc("/fivempolicelog", fivempolicelog.HandleRequest)
	//	log.Printf("Registered FiveM police log service on /fivempolicelog")
	//}

}
