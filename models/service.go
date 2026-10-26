package models

import "time"

// ServiceHistory merepresentasikan catatan riwayat pengerjaan/servis kendaraan konsumen.
type ServiceHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CustomerID      uint      `gorm:"not null;index" json:"customer_id"`
	ServiceDate     string    `gorm:"not null" json:"service_date"` // Format: YYYY-MM-DD
	Mileage         int       `json:"mileage"`                      // Kilometer saat servis
	Description     string    `gorm:"not null" json:"description"`  // Rincian servis (ganti oli, rem, dll.)
	NextServiceDate string    `json:"next_service_date"`            // Estimasi jadwal servis berikutnya
	NextMileage     int       `json:"next_mileage"`                 // Estimasi KM servis berikutnya
	Notes           string    `json:"notes"`                        // Catatan teknis / keluhan
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel di database
func (ServiceHistory) TableName() string {
	return "gas_services"
}
