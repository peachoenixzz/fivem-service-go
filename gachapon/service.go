package gachapon

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
