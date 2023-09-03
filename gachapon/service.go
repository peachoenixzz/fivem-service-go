package gachapon

import (
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

func handleRandGachaponItems(gci []GachaponItem) (*Item, float64) {
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
			return &selectedItem.Item, selectedItem.PullRate
		}
	}

	return nil, 0.0
}
