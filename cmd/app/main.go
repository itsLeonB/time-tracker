package main

import (
	"log"
	"os"

	"github.com/itsLeonB/time-tracker/internal/delivery/http/server"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	loadErr := godotenv.Load(curDir + "/.env")
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir)
	}

	srv := server.SetupServer()
	srv.Serve()
}
