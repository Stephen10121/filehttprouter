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

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

			data2, err := os.Open(file)

			if err != nil {
				panic(err)
			}

			defer data2.Close()

			var buf2 bytes.Buffer
			io.Copy(&buf2, data2)
			asString2 := buf2.String()

			is, _ := exists(root + "/layout.html")
			if is {
				data, err := os.Open(root + "/layout.html")

				if err != nil {
					panic(err)
				}

				defer data.Close()

				var buf bytes.Buffer
				io.Copy(&buf, data)
				asString := buf.String()

				if strings.Contains(asString, "<slot />") {
					asString = strings.Replace(asString, "<slot />", asString2, 1)
					fmt.Fprint(w, asString)
				} else {
					fmt.Println("<slot /> doesnt exist")
					fmt.Fprint(w, asString2)
				}
			} else {
				fmt.Fprint(w, asString2)
			}
		})
	}

	if len(config.CustomRoutes) > 0 {
		for _, route := range config.CustomRoutes {
			http.HandleFunc(route.Endpoint, route.Handler)
		}
	}

	// Setting the static folder and the endpoint if the user specified one.
	if config.StaticDirectory.DirectoryPath != "" && config.StaticDirectory.EndpointPath != "" {
		fs := http.FileServer(http.Dir(config.StaticDirectory.DirectoryPath))
		http.Handle(config.StaticDirectory.EndpointPath, http.StripPrefix(config.StaticDirectory.EndpointPath, fs))
	}

	fmt.Println("Running server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
