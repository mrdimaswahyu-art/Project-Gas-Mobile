package main

import (
	"log"
	"os"

	"go_vibe/config"
	"go_vibe/handlers"
	"go_vibe/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Muat file .env jika tersedia
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan variabel environment sistem.")
	}

	// Hubungkan ke database PostgreSQL
	if _, err := config.ConnectDatabase(); err != nil {
		log.Fatalf("[FATAL] Gagal inisialisasi database: %v", err)
	}

	router := gin.Default()

	// Sajikan file statis frontend (HTML, CSS, JS)
	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/login", "./web/login.html")
	router.StaticFile("/register", "./web/register.html")

	// Endpoint API Publik (Autentikasi)
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// Endpoint API Terlindungi (Memerlukan Token JWT)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/auth/me", handlers.Me)

		// Data Konsumen & Riwayat Servis
		api.GET("/brands", handlers.GetCarBrands)
		api.GET("/customers", handlers.GetCustomers)
		api.GET("/customers/:id", handlers.GetCustomerByID)
		api.POST("/customers", handlers.CreateCustomer)
		api.PUT("/customers/:id", handlers.UpdateCustomer)
		api.DELETE("/customers/:id", handlers.DeleteCustomer)

		// Riwayat Servis Khusus Konsumen
		api.POST("/customers/:id/services", handlers.AddServiceHistory)
		api.DELETE("/services/:id", handlers.DeleteServiceHistory)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Server Garuda Auto Service berjalan di http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
