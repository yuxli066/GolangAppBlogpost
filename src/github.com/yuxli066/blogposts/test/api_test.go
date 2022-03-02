package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app"

	"github.com/gavv/httpexpect/v2"
)

// declare types for test package
type apiTestCases []struct {
	url     string
	result  interface{}
	queries interface{}
}

type errorTestCases []struct {
	url     string
	result  string
	queries query
}

type query map[string]string

// Error Test Cases
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
			"tags":      "history,tech",
			"sortBy":    "likes",
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

	s := httpexpect.New(t, server.URL)

	for _, v := range api_tcs {
		if v.queries != nil {
			s.GET(v.url).WithQueryObject(v.queries.(query)).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("posts").Value("posts").Array().Length().Equal(v.result)
		} else {
			s.GET(v.url).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("success").ValueEqual("success", true)
		}
	}

}
