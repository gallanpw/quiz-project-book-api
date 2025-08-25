package main

import (
	"quiz-project-book-api/config"
	"quiz-project-book-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi ke database
	config.ConnectDB()

	// Inisialisasi router Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Jalankan server
	r.Run(":8080")
}
