package utils

func InArray(a interface{}, list []interface{}) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
