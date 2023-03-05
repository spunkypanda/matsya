package main

import (
	"fmt"
	"matsya/src/api"
	"matsya/src/config"
	"matsya/src/daemon"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	envName := config.GetEnvironmentName("ENV")
	if envName == "" {
		envName = "development"
	}

	pathPrefix := ""
	if envName == "test" {
		pathPrefix = "../"
	}

	config.Initialize(envName, pathPrefix)

	mode := config.GetString("mode")
	if mode == "" {
		mode = "api"
	}

	if mode == "daemon" {
		stop := make(chan string)
		quit := make(chan os.Signal, 1)

		watchedSignals := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL}
		signal.Notify(quit, watchedSignals...)

		go daemon.LongRunningProcess(stop)

		for {
			select {
			case msg := <-stop:
				fmt.Println(msg)
				quit <- os.Interrupt
			case <-quit:
				fmt.Println("Daemon interrupted. Quit listening!")
				return
			}
		}
	}

	address := config.GetString("host.domain")
	api.Serve(address)
}
