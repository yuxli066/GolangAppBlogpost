package errors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app"

	"github.com/gavv/httpexpect/v2"
)

func TestErrorHandling(t *testing.T) {
	// http handler setup & initialize routes for application
	app := app.App{}
	app.Initialize()
	handler := app.GetHTTPHandler()

	// setup golang api test server
	server := httptest.NewServer(handler)
	defer server.Close()

	s := httpexpect.New(t, server.URL)
	s.GET("/api/ping").Expect().Status(http.StatusOK)
}
