package main

import (
	"log"

	"github.com/NikhilSharmaWe/lib/app"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("vars.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	application := app.NewApplication()
	mux := application.Router()

	log.Fatal(mux.Start(application.Addr))
}
