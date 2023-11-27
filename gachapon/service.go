package gachapon

import (
	"fmt"
	"math/rand"
	"time"
)

func handleGachaponPlayer(pi map[string]int, ag []AllGachapon) []ResponsePlayerGachapon {
	var pgs []ResponsePlayerGachapon
	for _, item := range ag {
		if pi[item.Name] > 0 {
			pg := ResponsePlayerGachapon{
				Name:      item.Name,
				LabelName: item.LabelName,
				Quantity:  pi[item.Name],
			}
			pgs = append(pgs, pg)
			//fmt.Println(item.Name, " : ", pi[item.Name])
		}
	}
	return pgs
}

func handleRandResponseAndInsertGachapon(item map[string]map[string]any, gci []GachaponItem) ([]ItemInsert, []Item) {
	fmt.Println("Gachapon Summary:")
	var itemsInsert []ItemInsert
	var items []Item
	for _, gi := range gci {
		for itemId, i := range item {
			if itemId == gi.Item.ItemId {
				summary := fmt.Sprintf("%v (ID : %v): Player Got Count : %v , Total multipy qauntity : %v (Quantity Item per count : %v , PullRate : %v)",
					gi.Item.Name, itemId, i["count"], i["count"].(int)*gi.Item.Quantity, gi.Item.Quantity, gi.PullRate)
				fmt.Println(summary)

				itemInsert := ItemInsert{
					Name:       gi.Item.Name,
					ItemId:     itemId,
					Quantity:   i["count"].(int) * gi.Item.Quantity,
					GachaponID: i["gachapon_id"].(int),
					Category:   i["category"].(string),
				}

				item := Item{
					Name:     gi.Item.Name,
					ItemId:   itemId,
					Quantity: i["count"].(int),
					Category: i["category"].(string),
				}

				itemsInsert = append(itemsInsert, itemInsert)
				items = append(items, item)
			}
		}
	}
	return itemsInsert, items
}

func handleRandGachaponItems(gci []GachaponItem) (*Item, int, string) {
	secureRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randVal := secureRand.Float64()
	//pullSum := 0.0
	var sameRateItems []GachaponItem

	for _, gi := range gci {
		if randVal < gi.PullRate {
			sameRateItems = append(sameRateItems, gi)
		}
	}

	if len(sameRateItems) > 0 {
		randIndex := rand.Intn(len(sameRateItems))
		selectedItem := sameRateItems[randIndex]
		if selectedItem.Item.Quantity > 0 {
			selectedItem.Item.Quantity--
			return &selectedItem.Item, selectedItem.GachaponID, selectedItem.Item.Category
		}
	}

	return nil, 0, ""
}
