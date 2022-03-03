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
		url:     "/api/ping",
		result:  nil,
		queries: nil,
	},
	{
		url:    "/api/posts",
		result: 28,
		queries: query{
			"tags": "tech",
		},
	},
	{
		url:    "/api/posts",
		result: 46,
		queries: query{
			"tags": "history,tech",
		},
	},
	{
		url:    "/api/posts",
		result: 68,
		queries: query{
			"tags": "tech,health,history",
		},
	},
	{
		url:    "/api/posts",
		result: 26,
		queries: query{
			"tags": "history",
		},
	},
	{
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "id",
		},
	},
	{
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "reads",
		},
	},
	{
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "likes",
		},
	},
	{
		url:    "/api/posts",
		result: "asc",
		queries: query{
			"tags":   "history",
			"sortBy": "popularity",
		},
	},
	{
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "id",
			"direction": "desc",
		},
	},
	{
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "reads",
			"direction": "desc",
		},
	},
	{
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "likes",
			"direction": "desc",
		},
	},
	{
		url:    "/api/posts",
		result: "desc",
		queries: query{
			"tags":      "history",
			"sortBy":    "popularity",
			"direction": "desc",
		},
	},
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
		if v.queries != nil {
			switch v.result.(type) {
			case string:
				res := s.GET(v.url).WithQueryObject(v.queries.(query)).Expect().Status(http.StatusOK).JSON().Object().Value("posts").Array()
				var actualSortedValues []float64
				var expectedSortedValues []float64
				for _, field := range res.Iter() {
					var n float64
					n = field.Object().Value(v.queries.(query)["sortBy"]).Number().Raw()
					actualSortedValues = append(actualSortedValues, n)
					expectedSortedValues = append(expectedSortedValues, n)
				}
				if v.result.(string) == "asc" {
					sort.Slice(expectedSortedValues, func(i, j int) bool { return expectedSortedValues[i] < expectedSortedValues[j] })
				} else {
					sort.Slice(expectedSortedValues, func(i, j int) bool { return expectedSortedValues[i] > expectedSortedValues[j] })
				}
				assert.Equal(t, expectedSortedValues, actualSortedValues)
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
