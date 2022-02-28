package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

// use wait groups here to run concurrent requests to hatchways api
var m = sync.RWMutex{}

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(100)
	tags := r.URL.Query()["tags"]
	if tags == nil {
		respondError(w, http.StatusBadRequest, "Missing query parameter 'tags'")
	} else {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "https://api.hatchways.io/assessment/blog/posts", nil)
		if err != nil {
			log.Fatal(err)
		}

		tagQueries := req.URL.Query()
		wg := sync.WaitGroup{}
		strReceiver := make(chan []byte)

		// use string builder for response
		var sb strings.Builder
		for i := 0; i < len(tags); i++ {
			wg.Add(1)
			tagQueries.Add("tag", tags[i])
			req.URL.RawQuery = tagQueries.Encode()
			fmt.Println(req.URL)
			m.Lock()
			go getPostData(client, req, strReceiver, &wg)
			tagQueries.Del("tag")
			sb.WriteString(string(<-strReceiver))
		}
		// if query contains sortBy parameter, do something
		// if query contains direction parameter, do something
		respondJSON(w, http.StatusOK, sb)

	}
}

func getPostData(client *http.Client, request *http.Request, receiver chan<- []byte, wg *sync.WaitGroup) {
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	m.Unlock()
	wg.Done()

	fmt.Println(body)
	receiver <- body
}
