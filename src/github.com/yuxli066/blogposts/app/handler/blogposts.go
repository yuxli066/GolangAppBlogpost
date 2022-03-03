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
	"time"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app/utils"

	"github.com/patrickmn/go-cache"
)

// instantiate go cache
var Cache = cache.New(5*time.Minute, 5*time.Minute)

// use wait groups here to run concurrent requests to hatchways api
var m = sync.RWMutex{}

// constant string slices
const hatchwaysAPI string = "https://api.hatchways.io/assessment/blog/posts"

func getSortByFields() []string {
	return []string{"id", "reads", "likes", "popularity"}
}
func getDirectionFields() []string {
	return []string{"asc", "desc"}
}

// API Handler functions
func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(100)
	// default sort by & sort direction values
	var sortByField string = "id"
	var sortDirectionField string = "asc"

	queryTags := r.URL.Query()["tags"]
	querySortBy := r.URL.Query()["sortBy"]
	querySortDirection := r.URL.Query()["direction"]

	if queryTags == nil {
		respondError(w, http.StatusBadRequest, "Tags parameter is required")
	} else {
		tags := strings.Split(queryTags[0], ",")
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, hatchwaysAPI, nil)
		if err != nil {
			log.Fatal(err)
			respondError(w, http.StatusInternalServerError, "Hatchaway API Error")
		}

		// get sort by field & sort direction field
		if querySortBy != nil {
			if !utils.SliceContains(getSortByFields(), querySortBy[0]) {
				respondError(w, http.StatusBadRequest, "The sortBy parameter must be id, reads, likes, or popularity")
				return
			}
			sortByField = querySortBy[0]
		}

		if querySortDirection != nil {
			if !utils.SliceContains(getDirectionFields(), querySortDirection[0]) {
				respondError(w, http.StatusBadRequest, "The sortBy direction must be either asc or desc")
				return
			}
			sortDirectionField = querySortDirection[0]
		}

		// simple golang cache storage
		cacheKey := "tags-" + queryTags[0] + "sortBy-" + sortByField + "direction-" + sortDirectionField

		// respond with data from cache if found, otherwise save result to cache and respond
		d, found := Cache.Get(cacheKey)
		if found {
			respondJSON(w, http.StatusOK, d)
		} else {
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
				go utils.MergeMaps(mergedDataMap, dataMap, &wg, &m)
			}
			wg.Wait() // wait for goroutines to finish executing

			// sort results based on parameters
			sort.Slice(mergedDataMap["posts"], utils.CustomSort(mergedDataMap["posts"].([]interface{}), sortByField, sortDirectionField))
			Cache.Set(cacheKey, mergedDataMap, cache.DefaultExpiration)

			respondJSON(w, http.StatusOK, mergedDataMap)
		}
	}
}

// Makes API calls to hatchways api for blogposts payload
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
