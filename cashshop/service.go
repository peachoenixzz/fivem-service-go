package cashshop

func HandleLimitType(res ResponseValidateItem) bool {
	if res.LimitType == "01" || res.LimitType == "02" {
		if res.RemainQuantity == 0 {
			return false
		}
	}
	return true
}

func HandleMessage(count int64) Message {
	if count > 0 {
		return Message{"success"}
	}
	return Message{"fail"}
}

func HandleDateExpire(ex int) bool {
	if ex > 0 {
		return true
	}
	return false
}
