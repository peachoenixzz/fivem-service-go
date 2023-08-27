package playerquest

import (
	"fmt"
	"math/rand"
	"time"
)

func handleQuestItem(res []ResponseQuestItem) []ResponseSelectedItem {
	secureRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Shuffle the items to get random selection
	secureRand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	// Select 3-5 items
	selectedCount := secureRand.Intn(3) + 3
	selectedItems := res[:selectedCount]

	var rsis []ResponseSelectedItem
	var rsi ResponseSelectedItem
	// Assign random quantities based on rarity and populate selectedResponseItems
	for _, item := range selectedItems {
		switch item.Rare {
		case "normal":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(7) + 4)
		case "medium":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(3) + 1)
		case "hard":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(2) + 1)
		}

		rsis = append(rsis, rsi)
	}

	//for _, item := range rsis {
	//	fmt.Println("Item Name:", item.Name)
	//	fmt.Println("Item Rareness:", item.Rare)
	//	fmt.Println("Item Quantity:", item.Quantity)
	//}

	return rsis
}

func handleComparePlayerAndQuestItem(pi []PlayerItems, rqpi []ResponsePlayerQuestItem) {

	for _, items := range rqpi {
		fmt.Println(pi[items.Quantity], " / ", items.Quantity)
	}

}
