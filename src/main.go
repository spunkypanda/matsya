package main

import (
	"matsya/src/api"
	"matsya/src/config"

	_ "github.com/lib/pq"
)

func main() {
	envName := "development"
	config.Initialize(envName, "")

	address := config.GetString("host.domain")
	api.Serve(address)
}
