package main

import (
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	defaultConfigs := server.DefaultConfigs()
	ezutil.RunServer(defaultConfigs, server.SetupHTTPServer)
}
