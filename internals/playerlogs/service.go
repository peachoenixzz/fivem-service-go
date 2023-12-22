package playerlogs

import (
	"fmt"
	contime "github.com/kkgo-software-engineering/workshop/pkg/converttime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func selectQuery(req RequestCustomLog) bson.M {
	beginTime := contime.ParseTime(req.Begin)
	untilTime := contime.ParseTime(req.Until)

	plusGMTBegin := beginTime.In(time.FixedZone("GMT+7", 7*60*60))
	plusGMTUntil := untilTime.In(time.FixedZone("GMT+7", 7*60*60))
	if req.DiscordID == "" && req.Event == "" {
		f := bson.M{
			"timestamp": bson.M{
				"$gte": primitive.NewDateTimeFromTime(plusGMTBegin),
				"$lte": primitive.NewDateTimeFromTime(plusGMTUntil),
			},
			"content": bson.M{
				"$regex": req.Regex,
			},
		}
		return f
	}

	if req.DiscordID != "" && req.Event != "" {
		f := bson.M{
			"player.identifiers.discord": fmt.Sprintf("discord:%v", req.DiscordID),
			"event":                      req.Event,
			"timestamp": bson.M{
				"$gte": primitive.NewDateTimeFromTime(plusGMTBegin),
				"$lte": primitive.NewDateTimeFromTime(plusGMTUntil),
			},
			"content": bson.M{
				"$regex": req.Regex,
			},
		}
		return f
	}

	if req.DiscordID != "" && req.Event == "" {
		f := bson.M{
			"player.identifiers.discord": fmt.Sprintf("discord:%v", req.DiscordID),
			"timestamp": bson.M{
				"$gte": primitive.NewDateTimeFromTime(plusGMTBegin),
				"$lte": primitive.NewDateTimeFromTime(plusGMTUntil),
			},
			"content": bson.M{
				"$regex": req.Regex,
			},
		}
		return f
	}

	if req.DiscordID == "" && req.Event != "" {
		f := bson.M{
			"event": req.Event,
			"timestamp": bson.M{
				"$gte": primitive.NewDateTimeFromTime(plusGMTBegin),
				"$lte": primitive.NewDateTimeFromTime(plusGMTUntil),
			},
			"content": bson.M{
				"$regex": req.Regex,
			},
		}
		return f
	}
	return nil
}
