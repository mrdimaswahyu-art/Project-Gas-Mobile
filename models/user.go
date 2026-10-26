package models

import "time"

// User merepresentasikan staf atau pengguna bengkel yang dapat login ke sistem.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // Password hash, tidak diekspos di JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel di database agar tidak bertabrakan dengan tabel latihan lain
func (User) TableName() string {
	return "gas_users"
}
