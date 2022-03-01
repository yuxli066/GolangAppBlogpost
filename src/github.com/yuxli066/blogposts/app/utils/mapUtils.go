package utils

import (
	"sync"
)

func MergeMaps(mergedMap map[string]interface{}, newMap map[string]interface{}, wg *sync.WaitGroup, m *sync.RWMutex) {
	if len(mergedMap) == 0 {
		for k, v := range newMap { // deep copy new map to merged map
			mergedMap[k] = v
		}
	} else {
		newMapList := newMap["posts"].([]interface{})
		for _, b := range newMapList {
			switch t := b.(type) {
			case map[string]interface{}:
				if !Contains(mergedMap["posts"].([]interface{}), t["id"]) {
					mergedMap["posts"] = append(mergedMap["posts"].([]interface{}), b)
				}
				break
			}
		}
	}
	m.Unlock()
	wg.Done()
}
