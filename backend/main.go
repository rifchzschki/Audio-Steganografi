package main

import (
	"log"

	"github.com/rifchzschki/Audio-Steganografi/backend/routes"
)

func main() {
	router := routes.SetupRouter()

	log.Fatal(router.Run(":8080")) 
}