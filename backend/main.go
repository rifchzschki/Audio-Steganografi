package main

import (
	"log"
	"os"

	"github.com/rifchzschki/Audio-Steganografi/backend/cli"
	"github.com/rifchzschki/Audio-Steganografi/backend/routes"
)

func main() {

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "cli" {
			// Jalankan mode CLI
			cli.Run(os.Args[2:])
			return
		}
	}
	router := routes.SetupRouter()

	log.Fatal(router.Run(":8080")) 
}