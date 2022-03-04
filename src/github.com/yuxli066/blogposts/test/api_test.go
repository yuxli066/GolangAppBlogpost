package test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app"
	"leo-blog-post/src/github.com/yuxli066/blogposts/app/utils"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

// declare types for test package
type apiTestCases []struct {
	testName string
	url      string
	result   interface{}
	queries  interface{}
}

type errorTestCases []struct {
	testName string
	url      string
	result   string
	queries  query
}

type query map[string]string

// API Test Cases
var api_tcs apiTestCases = apiTestCases{
	{
		testName: "test api health check",
		url:     "/api/ping",
		result:  nil,
		queries: nil,
	},
	{
		testName: "test api tags post",
		url:    "/api/posts",
		result: 28,
		queries: query{
			"tags": "tech",
		},
	},
	{
		testName: "test api tags history, tech",
		url:    "/api/posts",
		result: 46,
		queries: query{
			"tags": "history,tech",
		},
	},
	{
		testName: "test api tags tech, health, history",
		url:    "/api/posts",
		result: 68,
		queries: query{
			"tags": "tech,health,history",
		},
	},
	{
		testName: "test api tags history",
		url:    "/api/posts",
		result: 26,
		queries: query{
			"tags": "history",
		},
	},
	{
		testName: "test api tags history sort by id in ascending order",
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "id",
		},
	},
	{
		testName: "test api tags history sort by reads in ascending order",
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "reads",
		},
	},
	{
		testName: "test api tags history sort by likes in ascending order",
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "likes",
		},
	},
	{
		testName: "test api tags history sort by popularity in ascending order",
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "popularity",
		},
	},
	{
		testName: "test api tags history sort by id in descending order",
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "id",
			"direction": "desc",
		},
	},
	{
		testName: "test api tags history sort by reads in descending order",
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "reads",
			"direction": "desc",
		},
	},
	{
		testName: "test api tags history sort by likes in descending order",
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "likes",
			"direction": "desc",
		},
	},
	{
		testName: "test api tags history sort by popularity in descending order",
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "popularity",
			"direction": "desc",
		},
	},
	{
		testName: "test api tags tech return correct results (correct tags)",
		url:    "/api/posts",
		result: []string{"tech","history","startups"},
		queries: query{
			"tags":      "tech,history,startups",
		},
	},
	// {
	// 	testName: "test api tags tech,history,startups return correct results (correct tags)",
	// 	url:    "/api/posts",
	// 	result: []string{"history"},
	// 	queries: query{
	// 		"tags":      "history",
	// 	},
	// },
}

func TestAPIFunctionality(t *testing.T) {
	// http handler setup & initialize routes for application
	app := app.App{}
	app.Initialize()
	handler := app.GetHTTPHandler()

	// setup golang api test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// setup console logging
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	s := httpexpect.New(t, server.URL)

	for _, v := range api_tcs {
		// log.Println(v.testName)
		if v.queries != nil {
			switch ty := v.result.(type) {
			case string, []string:
				_, isString := ty.(string)
				switch isString {
				case true: 
					res := s.GET(v.url).WithQueryObject(v.queries.(query)).Expect().Status(http.StatusOK).JSON().Object().Value("posts").Array()
					var actualSortedValues []float64
					var expectedSortedValues []float64
					for _, field := range res.Iter() {
						var n float64
						n = field.Object().Value(v.queries.(query)["sortBy"]).Number().Raw()
						actualSortedValues = append(actualSortedValues, n)
						expectedSortedValues = append(expectedSortedValues, n)
					}
					if ty == "asc" {
						sort.Slice(expectedSortedValues, func(i, j int) bool { return expectedSortedValues[i] < expectedSortedValues[j] })
					} else {
						sort.Slice(expectedSortedValues, func(i, j int) bool { return expectedSortedValues[i] > expectedSortedValues[j] })
					}
					assert.Equal(t, expectedSortedValues, actualSortedValues)
					break 
				case false: 
					res := s.GET(v.url).WithQueryObject(v.queries.(query)).Expect().Status(http.StatusOK).JSON().Object().Value("posts").Array()
					for _, field := range res.Iter() {
						containsTag := false
						for _, tag := range field.Object().Value("tags").Array().Iter() {
							tagString := tag.String().Raw()
							if utils.SliceContains(v.result.([]string), tagString) {
								containsTag = true
							}
						}
						assert.Equal(t, containsTag, true)
					}
					break
				default: 
					break
				}
				break
			case int:
				s.GET(v.url).WithQueryObject(v.queries.(query)).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("posts").Value("posts").Array().Length().Equal(v.result)
				break
			default:
				break
			}
		} else {
			s.GET(v.url).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("success").ValueEqual("success", true)
		}
	}
	t.Log(buf.String())
}
