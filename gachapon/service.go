package gachapon

func handleGachaponPlayer(pi map[string]int, ag []AllGachapon) []ResponsePlayerGachapon {
	var pgs []ResponsePlayerGachapon
	for _, item := range ag {
		pg := ResponsePlayerGachapon{
			Name:      item.Name,
			LabelName: item.LabelName,
			Quantity:  pi[item.Name],
		}
		pgs = append(pgs, pg)
		//fmt.Println(items.ItemName, " : ", pi[items.ItemName], " / ", items.Quantity)
	}
	return pgs
}
