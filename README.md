# File Http Router
This is a simple http web server that uses folder based routing. For example, lets say you have a directory called "app" and inside the directory, the structure looks like this:
- app
    - index.html
    - about
        - index.html

My framework will convert the directory structure into http routes. Like these:

/ -> will serve /app/index.html file

/about -> will serve the /app/about/index.html file
# Getting Started
First create an empty folder or repo. Then initailize the go package.
```bash
go mod init
```
Install my package.
```bash
go get github.com/stephen10121/filehttprouter@latest
```
Create a `main.go` file and the only thing you need in there to get started is this:
```go
package main

import "github.com/stephen10121/filehttprouter"

func main() {
	app := filehttprouter.App{}

	app.Run()
}
```
Create a `app` folder and inside is where you can create you index.html files and nested routes. Like this:
- app
    - index.html
    - about
        - index.html
    - services
        - index.html
        - nestedpage1
            - index.html
        - nestedpage2
            - index.html

To run the app, the command is
```bash
go run main.go
```
To build the app, the command is
```bash
go build main.go
```
running this command creates an executable that you can run.

# Features
These are the features that are currently working and in production.
## Custom Routes
If you need custom routes that returns json, file-streaming, plain text, etc, you can setup custom routes easily.
```go
func CustomEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := []int{2, 4}

	json.NewEncoder(w).Encode(p)
}

func main() {
	app := filehttprouter.App{
		CustomRoutes:    []filehttprouter.CustomRoute{
            {
                Endpoint: "/custom",
                Handler:  CustomEndpoint,
            },
        },
	}

	app.Run()
}
```
## Custom app directory path
If you dont want ./app to be the path of the route based directory, simply change the configuration.
```go
func main() {
	app := filehttprouter.App{
        Path:            "./anotherdirectorypath",
	}

	app.Run()
}
```
## Custom Port
If you want a different port, simply change the configuration.
```go
func main() {
	app := filehttprouter.App{
		Port:            "9000",
	}

	app.Run()
}
```
## Static Folder
You can also setup a static folder for your images, stylesheets, etc. This is how to setup the static folder and the api endpoint for it.

```go
func main() {
	app := filehttprouter.App{
        StaticDirectory: filehttprouter.StaticPath{
            DirectoryPath: "./static", // you can choose the static folder
            EndpointPath:  "/static/", // as well as its api endpoint
        },
	}

	app.Run()
}
```
# Example App
This is how an example app would look like.
```go
package main

import (
	"encoding/json"
	"net/http"

    // If you're tired of using "filehttprouter" everywhere, just put "router" (or any other variable) in front if the import string.
	router "github.com/stephen10121/filehttprouter"
)

func CustomEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := []int{2, 4}

	json.NewEncoder(w).Encode(p)
}

var myCustomRoutes = []router.CustomRoute{
	{
		Endpoint: "/custom",
		Handler:  CustomEndpoint,
	},
}

var myStaticDirectory = router.StaticPath{
	DirectoryPath: "./static",
	EndpointPath:  "/static/",
}

func main() {
	app := router.App{
		Path:            "./routes",
		Port:            "9000",
		CustomRoutes:    myCustomRoutes,
		StaticDirectory: myStaticDirectory,
	}

	app.Run()
}
```

# For me
to update tag:
```
git tag v0.1.0
```
to push tag:
```
git push origion v0.1.0
```