package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go_vibe/config"
	"go_vibe/middleware"
	"go_vibe/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", Register)
		authGroup.POST("/login", Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/auth/me", Me)
		api.GET("/brands", GetCarBrands)
		api.GET("/stats/monthly", GetMonthlyStats)
		api.GET("/customers", GetCustomers)
		api.GET("/customers/:id", GetCustomerByID)
		api.POST("/customers", CreateCustomer)
		api.PUT("/customers/:id", UpdateCustomer)
		api.DELETE("/customers/:id", DeleteCustomer)
		api.POST("/customers/:id/services", AddServiceHistory)
		api.DELETE("/services/:id", DeleteServiceHistory)
	}

	return r
}

func initDB(t *testing.T) {
	_ = godotenv.Load("../.env")
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL tidak disetel")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests-12345")
	}
	_, err := config.ConnectDatabase()
	if err != nil {
		t.Fatalf("Gagal inisialisasi DB: %v", err)
	}
}

func TestAuthAndCustomerFlow(t *testing.T) {
	initDB(t)
	router := setupTestRouter()

	// 1. Register User Baru
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000)
	username := fmt.Sprintf("staf_%d", rnd)
	password := "password123"

	regPayload, _ := json.Marshal(AuthRequest{
		Username: username,
		Password: password,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(regPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Register gagal: status %d, body %s", w.Code, w.Body.String())
	}

	// 2. Login User
	loginPayload, _ := json.Marshal(AuthRequest{
		Username: username,
		Password: password,
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login gagal: status %d, body %s", w.Code, w.Body.String())
	}

	var loginResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token JWT tidak diterima saat login")
	}

	// 3. Tes Akses Endpoint Terlindungi (Create Customer)
	custPayload, _ := json.Marshal(CustomerRequest{
		Name:        "Budi Santoso",
		Phone:       "081234567890",
		Address:     "Jl. Merdeka No. 10",
		PlateNumber: "B 1234 XYZ",
		CarBrand:    "Toyota",
		CarType:     "Avanza 1.3 G",
	})

	req, _ = http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(custPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create customer gagal: status %d, body %s", w.Code, w.Body.String())
	}

	var createdCust models.Customer
	_ = json.Unmarshal(w.Body.Bytes(), &createdCust)
	if createdCust.ID == 0 {
		t.Fatal("Customer ID yang dibuat harus > 0")
	}

	// 4. Tambah Riwayat Servis
	servicePayload, _ := json.Marshal(ServiceRequest{
		ServiceDate:     "2026-10-26",
		Mileage:         15000,
		Description:     "Ganti oli mesin, filter oli, dan pengecekan rem",
		NextServiceDate: "2027-04-26",
		NextMileage:     20000,
		Notes:           "Kondisi kampas rem masih 80%",
	})

	serviceURL := fmt.Sprintf("/api/customers/%d/services", createdCust.ID)
	req, _ = http.NewRequest(http.MethodPost, serviceURL, bytes.NewBuffer(servicePayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Tambah servis gagal: status %d, body %s", w.Code, w.Body.String())
	}

	// 5. Statistik Bulanan (12 bulan terakhir)
	req, _ = http.NewRequest(http.MethodGet, "/api/stats/monthly?months=12", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Statistik bulanan gagal: status %d, body %s", w.Code, w.Body.String())
	}

	var stats MonthlyStatsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Respons statistik bulanan tidak valid: %v", err)
	}
	if len(stats.Months) != 12 {
		t.Fatalf("Jumlah bulan statistik harus 12, didapat %d", len(stats.Months))
	}
	if stats.Totals.NewCustomers < 1 {
		t.Fatal("Total konsumen baru harus minimal 1 setelah membuat konsumen")
	}
	if stats.Totals.Services < 1 {
		t.Fatal("Total servis harus minimal 1 setelah menambah riwayat servis")
	}

	// 6. Pencarian Pelanggan berdasarkan Plat Nomor
	searchURL := "/api/customers?search=1234"
	req, _ = http.NewRequest(http.MethodGet, searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Pencarian customer gagal: status %d", w.Code)
	}

	var searchResults []models.Customer
	_ = json.Unmarshal(w.Body.Bytes(), &searchResults)
	if len(searchResults) == 0 {
		t.Fatal("Hasil pencarian customer kosong, harusnya ditemukan")
	}

	// 7. Filter Pelanggan berdasarkan Merk
	filterURL := "/api/customers?brand=Toyota"
	req, _ = http.NewRequest(http.MethodGet, filterURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Filter brand gagal: status %d", w.Code)
	}

	// 8. Cleanup Data Test
	config.DB.Delete(&createdCust)
}
