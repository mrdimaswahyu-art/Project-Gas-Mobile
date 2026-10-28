package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go_vibe/config"

	"github.com/gin-gonic/gin"
)

// monthlyCountRow menampung hasil agregasi COUNT per bulan dari database.
type monthlyCountRow struct {
	Month string
	Total int64
}

// MonthlyStatPoint adalah satu titik data pada diagram bulanan.
type MonthlyStatPoint struct {
	Month        string `json:"month"`         // Format: YYYY-MM
	Label        string `json:"label"`         // Contoh: Okt
	Year         int    `json:"year"`          // Contoh: 2026
	NewCustomers int64  `json:"new_customers"` // Jumlah konsumen baru pada bulan tersebut
	Services     int64  `json:"services"`      // Jumlah riwayat servis pada bulan tersebut
}

// MonthlyStatsTotals adalah total keseluruhan pada rentang bulan yang diminta.
type MonthlyStatsTotals struct {
	NewCustomers int64 `json:"new_customers"`
	Services     int64 `json:"services"`
}

// MonthlyStatsResponse adalah bentuk respons API /api/stats/monthly.
type MonthlyStatsResponse struct {
	StartMonth string             `json:"start_month"`
	EndMonth   string             `json:"end_month"`
	Months     []MonthlyStatPoint `json:"months"`
	Totals     MonthlyStatsTotals `json:"totals"`
}

var indonesianMonthNames = [...]string{
	"Jan", "Feb", "Mar", "Apr", "Mei", "Jun",
	"Jul", "Agu", "Sep", "Okt", "Nov", "Des",
}

// GetMonthlyStats mengembalikan jumlah konsumen baru dan riwayat servis per bulan.
// Parameter opsional: ?months=12 (minimal 1, maksimal 24, default 12).
func GetMonthlyStats(c *gin.Context) {
	months := 12
	if rawMonths := c.Query("months"); rawMonths != "" {
		parsed, err := strconv.Atoi(rawMonths)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter months harus berupa angka antara 1 sampai 24"})
			return
		}
		if parsed > 24 {
			parsed = 24
		}
		months = parsed
	}

	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startMonth := currentMonth.AddDate(0, -(months - 1), 0)

	var customerRows []monthlyCountRow
	if err := config.DB.Raw(`
		SELECT TO_CHAR(created_at, 'YYYY-MM') AS month, COUNT(*) AS total
		FROM gas_customers
		WHERE created_at >= ?
		GROUP BY 1`, startMonth).Scan(&customerRows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil statistik konsumen"})
		return
	}

	var serviceRows []monthlyCountRow
	if err := config.DB.Raw(`
		SELECT SUBSTRING(service_date FROM 1 FOR 7) AS month, COUNT(*) AS total
		FROM gas_services
		WHERE service_date >= ?
		GROUP BY 1`, startMonth.Format("2006-01-02")).Scan(&serviceRows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil statistik servis"})
		return
	}

	customerByMonth := make(map[string]int64, len(customerRows))
	for _, row := range customerRows {
		customerByMonth[row.Month] = row.Total
	}

	serviceByMonth := make(map[string]int64, len(serviceRows))
	for _, row := range serviceRows {
		serviceByMonth[row.Month] = row.Total
	}

	// Susun data berurutan dari bulan tertua ke bulan terbaru,
	// termasuk bulan yang belum memiliki data (nilainya 0).
	points := make([]MonthlyStatPoint, 0, months)
	totals := MonthlyStatsTotals{}
	for i := 0; i < months; i++ {
		month := startMonth.AddDate(0, i, 0)
		key := month.Format("2006-01")

		point := MonthlyStatPoint{
			Month:        key,
			Label:        indonesianMonthNames[int(month.Month())-1],
			Year:         month.Year(),
			NewCustomers: customerByMonth[key],
			Services:     serviceByMonth[key],
		}

		totals.NewCustomers += point.NewCustomers
		totals.Services += point.Services
		points = append(points, point)
	}

	c.JSON(http.StatusOK, MonthlyStatsResponse{
		StartMonth: startMonth.Format("2006-01"),
		EndMonth:   currentMonth.Format("2006-01"),
		Months:     points,
		Totals:     totals,
	})
}
