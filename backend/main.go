package main

import "github.com/johnwongx/webook/backend/integration"

func main() {
	server := integration.InitWebServer()

	server.Run(":8080")
}
