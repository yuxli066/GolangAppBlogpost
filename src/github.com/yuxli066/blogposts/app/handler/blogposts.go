package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"sync"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app/utils"
)

// use wait groups here to run concurrent requests to hatchways api
var m = sync.RWMutex{}

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(100)
	queryTags := r.URL.Query()["tags"]
	querySortBy := r.URL.Query()["sortBy"]
	querySortDirection := r.URL.Query()["direction"]

	if queryTags == nil {
		respondError(w, http.StatusBadRequest, "Tags parameter is required")
	} else {
		tags := strings.Split(queryTags[0], ",")
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "https://api.hatchways.io/assessment/blog/posts", nil)
		if err != nil {
			log.Fatal(err)
		}

		tagQueries := req.URL.Query()
		wg := sync.WaitGroup{}
		strReceiver := make(chan []byte)

		var data []byte // holds the return data from api call
		mergedDataMap := make(map[string]interface{})

		// Using wait groups and Mutexes to create concurrent http requests & data parsing
		for _, t := range tags {
			wg.Add(2)
			m.Lock()
			go getPostData(client, req, strReceiver, &wg, &tagQueries, t)
			data = <-strReceiver
			dataMap := make(map[string]interface{})
			json.Unmarshal(data, &dataMap)
			m.Lock()
			go mergeMaps(mergedDataMap, dataMap, &wg)
		}
		wg.Wait() // wait for goroutines to finish executing

		acceptedSortByFields := []string{"id", "reads", "likes", "popularity"}
		acceptedDirectionFields := []string{"asc", "desc"}

		// default sort by & sort direction values
		sortByField := "id"
		sortDirectionField := "asc"

		// get sort by field & sort direction field
		if querySortBy != nil {
			if !utils.SliceContains(acceptedSortByFields, querySortBy[0]) {
				respondError(w, http.StatusBadRequest, "The sortBy parameter must be id, reads, likes, or popularity")
			}
			sortByField = querySortBy[0]
		}

		if querySortDirection != nil {
			if !utils.SliceContains(acceptedDirectionFields, querySortDirection[0]) {
				respondError(w, http.StatusBadRequest, "The sortBy direction must be either asc or desc")
			}
			sortDirectionField = querySortDirection[0]
		}
		// sort results based on parameters
		sort.Slice(mergedDataMap["posts"], customSort(mergedDataMap["posts"].([]interface{}), sortByField, sortDirectionField))

		// JSON response for api
		respondJSON(w, http.StatusOK, mergedDataMap)
	}
}

func getPostData(client *http.Client, request *http.Request, receiver chan<- []byte, wg *sync.WaitGroup, tagQueries *url.Values, tag string) {
	tagQueries.Add("tag", tag)
	request.URL.RawQuery = tagQueries.Encode()
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	defer tagQueries.Del("tag")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	m.Unlock()
	wg.Done()

	receiver <- body
}

func mergeMaps(mergedMap map[string]interface{}, newMap map[string]interface{}, wg *sync.WaitGroup) {
	if len(mergedMap) == 0 {
		for k, v := range newMap { // deep copy new map to merged map
			mergedMap[k] = v
		}
	} else {
		newMapList := newMap["posts"].([]interface{})
		for _, b := range newMapList {
			switch t := b.(type) {
			case map[string]interface{}:
				if !utils.Contains(mergedMap["posts"].([]interface{}), t["id"]) {
					mergedMap["posts"] = append(mergedMap["posts"].([]interface{}), b)
				}
				break
			}
		}
	}
	m.Unlock()
	wg.Done()
}

type allPosts []interface{}

func customSort(p allPosts, keyVal string, sortDirection string) func(int, int) bool {
	return func(i, j int) bool {
		if sortDirection == "asc" {
			return p[i].(map[string]interface{})[keyVal].(float64) < p[j].(map[string]interface{})[keyVal].(float64)
		} else {
			return p[i].(map[string]interface{})[keyVal].(float64) > p[j].(map[string]interface{})[keyVal].(float64)
		}
	}
}
