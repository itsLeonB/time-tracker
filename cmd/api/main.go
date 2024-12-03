package main

import "github.com/itsLeonB/time-tracker/internal/config"

func main() {
	app := config.SetupApp()
	app.Serve()
}
