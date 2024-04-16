package playerlogs

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
)

func selectQuery(req RequestCustomLog) (bson.M, error) {
	const layout = "2006-01-02T15:04"

	beginTime, errBegin := time.Parse(layout, req.Begin)
	if errBegin != nil {
		return nil, fmt.Errorf("invalid begin time format: %v", errBegin)
	}
	untilTime, errUntil := time.Parse(layout, req.Until)
	if errUntil != nil {
		return nil, fmt.Errorf("invalid until time format: %v", errUntil)
	}

	gmtPlus7 := time.FixedZone("GMT+7", 7*60*60)

	beginTimeGMT7 := beginTime.In(gmtPlus7)
	untilTimeGMT7 := untilTime.In(gmtPlus7)

	filter := bson.M{
		"timestamp": bson.M{
			"$gte": primitive.NewDateTimeFromTime(beginTimeGMT7),
			"$lte": primitive.NewDateTimeFromTime(untilTimeGMT7),
		},
		"content": bson.M{
			"$regex": req.Regex,
		},
	}

	// Conditionally add DiscordID and Event to the filter
	if req.DiscordID != "" {
		filter["player.identifiers.discord"] = fmt.Sprintf("discord:%v", req.DiscordID)
	}

	if req.Event != "" {
		filter["event"] = req.Event
	}

	// Return the constructed filter
	return filter, nil

}

func handleDiscordLog(req RequestInsert, h Handler) error {
	if req.Options.Important {
		_, err := h.Discord.ChannelMessageSend("1220302191013658644", fmt.Sprintf("Hello @everyone someone cheating your server now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "BuyCarFromDealerWithMoney" {
		_, err := h.Discord.ChannelMessageSend("1220383622217863208", fmt.Sprintf("Hello @everyone player buy vehicle now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "BuyCarFromDealerWithCredit" {
		_, err := h.Discord.ChannelMessageSend("1220383622217863208", fmt.Sprintf("Hello @everyone player buy vehicle now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "standardAddWeapon" {
		_, err := h.Discord.ChannelMessageSend("1220385374501605428", fmt.Sprintf("Hello @everyone player add weapon now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "ReceivedCraftItem" {
		_, err := h.Discord.ChannelMessageSend("1220385374501605428", fmt.Sprintf("Hello @everyone player craft item now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "ReceivedCraftItemFail" {
		_, err := h.Discord.ChannelMessageSend("1220385374501605428", fmt.Sprintf("Hello @everyone player craft item now ! \n %v", req.Content))
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}

	if req.Event == "standardAddInventoryItem" {
		item, count := parseItem(req.Content)
		if item != "" && isFilteredItem(item, []string{"cron", "weapon_box", "cement", "alloy", "exp", "heart_100", "gacha_support", "gacha_5", "gacha_event", "keycard_silver", "keycard_red", "keycard_gold", "ruby", "afk_gem", "gun_pin", "gun_spring", "gun_barrel", "heart_100", "keycard_silver", "keycard_gold", "keycard_red", "afk_squid"}) && count >= 7 {
			_, err := h.Discord.ChannelMessageSend("1220812542608146522", fmt.Sprintf("Hello @everyone player add item now ! \n `%v : ชื่อไอเทม : %v จำนวน : %v DiscordID : %v`", req.Player.Name, item, count, req.Player.Identifiers.Discord))
			if err != nil {
				return fmt.Errorf("unable to send message to Discord: %v", err)
			}
		}
	}
	return nil
}

func parseItem(input string) (string, int) {
	start := strings.Index(input, "(") + 1
	end := strings.Index(input, ")")
	parameters := strings.Split(input[start:end], ",")
	if len(parameters) != 2 {
		return "", 0
	}

	itemParts := strings.Split(parameters[0], "_")
	if len(itemParts) == 0 {
		return "", 0
	}

	count, err := strconv.Atoi(parameters[1])
	if err != nil {
		return "", 0
	}

	return itemParts[0], count
}

func isFilteredItem(item string, keywords []string) bool {
	for _, keyword := range keywords {
		if item == keyword {
			return true
		}
	}
	return false
}
