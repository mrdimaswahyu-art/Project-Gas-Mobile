package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestConnectDatabase(t *testing.T) {
	// Muat file .env dari root folder
	_ = godotenv.Load("../.env")

	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL tidak disetel, melewati pengujian database")
	}

	db, err := ConnectDatabase()
	if err != nil {
		t.Fatalf("ConnectDatabase() gagal: %v", err)
	}

	if db == nil {
		t.Fatal("Expected *gorm.DB tidak bernilai nil")
	}
}
