package handlers

import (
	"net/http"
	"strings"

	"go_vibe/config"
	"go_vibe/models"

	"github.com/gin-gonic/gin"
)

// CustomerRequest merepresentasikan payload input untuk membuat atau memperbarui data konsumen.
type CustomerRequest struct {
	Name        string `json:"name" binding:"required"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	PlateNumber string `json:"plate_number" binding:"required"`
	CarBrand    string `json:"car_brand" binding:"required"`
	CarType     string `json:"car_type" binding:"required"`
}

// ServiceRequest merepresentasikan payload input untuk menambah riwayat servis.
type ServiceRequest struct {
	ServiceDate     string `json:"service_date" binding:"required"`
	Mileage         int    `json:"mileage"`
	Description     string `json:"description" binding:"required"`
	NextServiceDate string `json:"next_service_date"`
	NextMileage     int    `json:"next_mileage"`
	Notes           string `json:"notes"`
}

// GetCustomers mengambil daftar konsumen dengan dukungan pencarian (plat nomor, tipe mobil, nama) dan filter merk mobil.
func GetCustomers(c *gin.Context) {
	search := strings.TrimSpace(c.Query("search"))
	brand := strings.TrimSpace(c.Query("brand"))

	query := config.DB.Model(&models.Customer{}).Preload("Services")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(plate_number) LIKE ? OR LOWER(car_type) LIKE ? OR LOWER(name) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if brand != "" {
		query = query.Where("LOWER(car_brand) = ?", strings.ToLower(brand))
	}

	var customers []models.Customer
	if err := query.Order("created_at desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data konsumen"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetCustomerByID mengambil detail 1 konsumen beserta riwayat servisnya.
func GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := config.DB.Preload("Services").First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data konsumen tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// CreateCustomer menambahkan data konsumen dan kendaraan baru.
func CreateCustomer(c *gin.Context) {
	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, plat nomor, merk mobil, dan tipe mobil wajib diisi"})
		return
	}

	customer := models.Customer{
		Name:        strings.TrimSpace(req.Name),
		Phone:       strings.TrimSpace(req.Phone),
		Address:     strings.TrimSpace(req.Address),
		PlateNumber: strings.ToUpper(strings.TrimSpace(req.PlateNumber)),
		CarBrand:    strings.TrimSpace(req.CarBrand),
		CarType:     strings.TrimSpace(req.CarType),
	}

	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data konsumen"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// UpdateCustomer memperbarui data konsumen dan kendaraannya.
func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data konsumen tidak ditemukan"})
		return
	}

	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	customer.Name = strings.TrimSpace(req.Name)
	customer.Phone = strings.TrimSpace(req.Phone)
	customer.Address = strings.TrimSpace(req.Address)
	customer.PlateNumber = strings.ToUpper(strings.TrimSpace(req.PlateNumber))
	customer.CarBrand = strings.TrimSpace(req.CarBrand)
	customer.CarType = strings.TrimSpace(req.CarType)

	if err := config.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data konsumen"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer menghapus data konsumen (otomatis menghapus riwayat servis karena CASCADE).
func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data konsumen tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data konsumen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data konsumen dan riwayat servis berhasil dihapus"})
}

// AddServiceHistory menambahkan riwayat servis ke kendaraan seorang konsumen.
func AddServiceHistory(c *gin.Context) {
	customerID := c.Param("id")
	var customer models.Customer

	if err := config.DB.First(&customer, customerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data konsumen tidak ditemukan"})
		return
	}

	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal servis dan deskripsi pengerjaan wajib diisi"})
		return
	}

	service := models.ServiceHistory{
		CustomerID:      customer.ID,
		ServiceDate:     strings.TrimSpace(req.ServiceDate),
		Mileage:         req.Mileage,
		Description:     strings.TrimSpace(req.Description),
		NextServiceDate: strings.TrimSpace(req.NextServiceDate),
		NextMileage:     req.NextMileage,
		Notes:           strings.TrimSpace(req.Notes),
	}

	if err := config.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan riwayat servis"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// DeleteServiceHistory menghapus satu riwayat servis.
func DeleteServiceHistory(c *gin.Context) {
	id := c.Param("id")
	var service models.ServiceHistory

	if err := config.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Riwayat servis tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat servis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Riwayat servis berhasil dihapus"})
}

// GetCarBrands mengembalikan daftar unik merk mobil yang sudah ada untuk opsi filter.
func GetCarBrands(c *gin.Context) {
	var brands []string
	if err := config.DB.Model(&models.Customer{}).Distinct("car_brand").Pluck("car_brand", &brands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar merk mobil"})
		return
	}

	c.JSON(http.StatusOK, brands)
}
