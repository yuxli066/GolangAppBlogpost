# GolangAppBlogpost
Golang Backend API for Hatchways Application

## Install 

* Run the command ```go mod download```
  - This command should download all the dependencies specified in the go.mod file 

## Build 

* Run the command ```go build ./src/github.com/yuxli066/blogposts/Main.go```
  - An executable file named "Main" should be generated in the parent/module directory 

## Running the App

* Run the command ```go run ./src/github.com/yuxli066/blogposts/Main.go```
* Run the command ```./Main```
    - If you have already built the executable already from the previous step

## Clearing Cache

* Run the command ```go clean -cache```

## Testing

* Run the command ```go test -v ./src/github.com/yuxli066/blogposts/test/```
    - If you wish to clear the cache before testing before running the tests run ```go clean -testcache``` first to clear cache

## File Structure
* The API source code can be found under the app directory ```go run ./src/github.com/yuxli066/blogposts/app/```
* The Main package can be found here ```go run ./src/github.com/yuxli066/blogposts/Main.go```
* The Test source code can be found under the app directory ```go run ./src/github.com/yuxli066/blogposts/test/```
* The Main executable file is included as well

## Credits

This software uses the following open source packages:

- [Golang](https://go.dev/)

This software used the following Golang Packages: 

- All packages can be found under the go.mod file

## License

MIT

---

> [yuxuanleoli.com](https://yuxuanleoli.com/) &nbsp;&middot;&nbsp;
> GitHub [@yuxli066]() &nbsp;&middot;&nbsp;

