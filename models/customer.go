package models

import "time"

// Customer merepresentasikan data pelanggan bengkel beserta informasi kendaraannya.
type Customer struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"not null" json:"name"`
	Phone       string           `json:"phone"`
	Address     string           `json:"address"`
	PlateNumber string           `gorm:"index;not null" json:"plate_number"` // Contoh: B 1234 ABC
	CarBrand    string           `gorm:"index;not null" json:"car_brand"`    // Contoh: Toyota, Honda, Suzuki
	CarType     string           `gorm:"index;not null" json:"car_type"`     // Contoh: Avanza, Brio, Ertiga
	Services    []ServiceHistory `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"services,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// TableName menentukan nama tabel di database
func (Customer) TableName() string {
	return "gas_customers"
}
