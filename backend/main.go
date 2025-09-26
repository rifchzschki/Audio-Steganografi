package main

import (
	"log"

	"github.com/rifchzschki/go-gin-react-project/backend/routes"
)

func main() {
	// Mendapatkan router yang sudah dikonfigurasi
	router := routes.SetupRouter()

	// Jalankan server pada port 8080
	log.Fatal(router.Run(":8080")) 
	// log.Fatal akan mencetak error jika server gagal dijalankan
}