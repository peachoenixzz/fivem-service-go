package cashshop

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkgo-software-engineering/workshop/config"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Cfg     config.FeatureFlag
	MongoDB *mongo.Client
	MysqlDB *sql.DB
}

func New(cfgFlag config.FeatureFlag, mongoDB *mongo.Client, mysqlDB *sql.DB) *Handler {
	return &Handler{cfgFlag, mongoDB, mysqlDB}
}

func (h Handler) getInitCashShop(ctx context.Context, p *mw.JwtCustomClaims) (ResponseInitCashShop, error) {

	return ResponseInitCashShop{}, nil
}
