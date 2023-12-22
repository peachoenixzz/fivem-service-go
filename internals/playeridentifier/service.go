package playeridentifier

func HandleTotalVipPoint(vp int64, ex int64, pp int64) int64 {
	tp := vp + ex + pp
	return tp
}

func HandleDateExpire(ex int) bool {
	if ex != 0 {
		return true
	}
	return false
}

//func HandlePlayerItem(req RequestUpdateVip) {
//	for _, item := range req.ExpireItems {
//		fmt.Println("Item Name : ", item.ItemName)
//		fmt.Println("Quantity : ", item.Quantity)
//		fmt.Println("Exp Date : ", item.ExpireDate)
//	}
//}
