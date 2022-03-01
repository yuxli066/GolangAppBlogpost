package utils

type allPosts []interface{}

func Contains(mergedMap allPosts, idNumber interface{}) bool {
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

func CustomSort(p allPosts, keyVal string, sortDirection string) func(int, int) bool {
	return func(i, j int) bool {
		if sortDirection == "asc" {
			return p[i].(map[string]interface{})[keyVal].(float64) < p[j].(map[string]interface{})[keyVal].(float64)
		} else {
			return p[i].(map[string]interface{})[keyVal].(float64) > p[j].(map[string]interface{})[keyVal].(float64)
		}
	}
}
