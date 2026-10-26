package config

import (
	"errors"
	"log"
	"os"

	"go_vibe/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase menginisialisasi koneksi PostgreSQL dan menjalankan migrasi tabel otomatis.
func ConnectDatabase() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("variabel lingkungan DATABASE_URL belum diisi")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, errors.New("gagal terhubung ke database: pastikan server database aktif dan URL benar")
	}

	// Migrasi otomatis tabel-tabel yang dibutuhkan
	err = db.AutoMigrate(&models.User{}, &models.Customer{}, &models.ServiceHistory{})
	if err != nil {
		return nil, errors.New("gagal melakukan migrasi tabel database")
	}

	DB = db
	log.Println("[INFO] Berhasil terhubung ke database dan migrasi tabel selesai.")
	return DB, nil
}
