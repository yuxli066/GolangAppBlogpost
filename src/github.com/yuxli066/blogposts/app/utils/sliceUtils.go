package utils

func Contains(mergedMap []interface{}, idNumber interface{}) bool {
	for _, s := range mergedMap {
		switch t := s.(type) {
		case map[string]interface{}:
			if t["id"].(float64) == idNumber.(float64) {
				return true
			}
			break
		}
	}
	return false
}

func SliceContains(sl []string, word string) bool {
	for _, v := range sl {
		if v == word {
			return true
		}
	}
	return false
}
