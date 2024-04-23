package main

import (
	test "github.com/stephen10121/iframe-test/iframetest"
)

func main() {
	app := test.App{
		Path: "./routes",
		Port: "9000",
	}

	app.Run()
}
