package main

import (
	"leo-blog-post/src/github.com/yuxli066/blogposts/app"
)

func main() {
	app := &app.App{}
	app.Run(":8080")
}