package filehttprouter

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stephen10121/filehttprouter/project"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

type CustomRoute struct {
	// This is where the api endpoint is. For example "/api/helloWorld". Warning if you have an endpoint that conflicts with the folder routed endpoint, this endpoint will overide the folder routed endpoint.
	Endpoint string
	// This is the function that handles the http request.
	Handler func(http.ResponseWriter, *http.Request)
}

type StaticPath struct {
	// The path of the static folder.
	DirectoryPath string
	// The api endpoint for the static files. Excample: /static/
	EndpointPath string
}

type App struct {
	// This is the path of the route based directory. The default path is ./app
	Path string
	// This sets the port of the server. Default port is 8080
	Port string
	// If you need a route that requires more than just an index.html file, use this.
	CustomRoutes []CustomRoute
	// This configures the static folder and its api endpoint
	StaticDirectory StaticPath
}

func (config App) Run() {
	root := "./app"
	if len(config.Path) > 0 {
		root = config.Path
	}

	port := "8080"
	if len(config.Port) > 0 {
		port = config.Port
	}

	files, err := FilePathWalkDir(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		dir, actualfile := filepath.Split(file)

		// For now, only index.html is suppoerted
		if actualfile != "index.html" {
			continue
		}

		dirToPath := ""

		splitDirectory := strings.Split(dir, project.PATH_SEPARATOR)
		for i := 1; i < len(splitDirectory); i++ {
			dirToPath += "/" + splitDirectory[i]
		}

		http.HandleFunc(dirToPath, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)

			data, err := os.Open(file)

			if err != nil {
				panic(err)
			}

			defer data.Close()

			fi, err := data.Stat()
			if err != nil {
				// Could not obtain stat, handle error
				panic(err)
			}

			w.Header().Set("Content-Length", fmt.Sprint(fi.Size()))

			var buf bytes.Buffer
			io.Copy(&buf, data)
			asString := buf.String()
			fmt.Fprint(w, asString)
		})
	}

	if len(config.CustomRoutes) > 0 {
		for _, route := range config.CustomRoutes {
			http.HandleFunc(route.Endpoint, route.Handler)
		}
	}

	// Setting the static folder and the endpoint if the user specified one.
	if config.StaticDirectory.DirectoryPath != "" && config.StaticDirectory.EndpointPath != "" {
		http.Handle(config.StaticDirectory.EndpointPath, http.FileServer(http.Dir(config.StaticDirectory.DirectoryPath)))
	}

	fmt.Println("Running server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
