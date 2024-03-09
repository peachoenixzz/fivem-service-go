package playerquest

import (
	"fmt"
	"math/rand"
	"time"
)

func handleCardAItem(res ResponseRequireQuestPlayer, pi map[string]int) ResponseRequireQuestPlayer {
	questAValue, found := pi["exp"]
	if found {
		res = ResponseRequireQuestPlayer{
			WeightLevel: res.WeightLevel,
			Quantity:    res.Quantity,
			CardAItem:   questAValue,
		}
	}

	if !found {
		res = ResponseRequireQuestPlayer{
			WeightLevel: res.WeightLevel,
			Quantity:    res.Quantity,
			CardAItem:   0,
		}
	}

	fmt.Println(res.CardAItem, res.WeightLevel, res.Quantity)
	return res
}

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
		fmt.Println(item.Rare)
		switch item.Rare {
		case "easy":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(6) + 3) // Generates a random number between 3 and 8 (inclusive)
		case "normal":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(4) + 2) // Generates a random number between 2 and 5 (inclusive)
		case "medium":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(3) + 1) // Generates a random number between 1 and 3 (inclusive)
		case "hard":
			rsi.Name = item.Name
			rsi.Quantity = int64(secureRand.Intn(2) + 1) // Generates a random number between 1 and 2 (inclusive)
		case "rare":
			rsi.Name = item.Name
			rsi.Quantity = 1 // Always set the quantity to 1 for the "rare" case
		}

		rsis = append(rsis, rsi)
	}

	for _, item := range rsis {
		fmt.Println("Item Name:", item.Name)
		fmt.Println("Item Quantity:", item.Quantity)
	}

	return rsis
}

func handleComparePlayerAndQuestItem(pi map[string]int, rqpi []ResponsePlayerQuestItem) []ResponseItemComparison {
	var comparisons []ResponseItemComparison
	for _, items := range rqpi {
		var comparsion ResponseItemComparison
		comparsion = ResponseItemComparison{
			ItemName:             items.ItemName,
			LabelName:            items.LabelName,
			Comparison:           fmt.Sprintf("%v/%v", pi[items.ItemName], items.Quantity),
			PlayerItemQuantity:   pi[items.ItemName],
			QuestRequireQuantity: items.Quantity,
		}
		comparisons = append(comparisons, comparsion)
		//fmt.Println(items.ItemName, " : ", pi[items.ItemName], " / ", items.Quantity)
	}
	return comparisons
}
