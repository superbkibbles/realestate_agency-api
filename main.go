package main

import (
	"github.com/joho/godotenv"
	"github.com/superbkibbles/realestate_agency-api/app"
)

func main() {
	godotenv.Load()
	app.StartApplication()
}
