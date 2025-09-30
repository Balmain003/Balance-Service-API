package stat

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"task/pkg/db"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(database *db.Db) *StatRepository {
	return &StatRepository{
		Database: database,
	}
}

func (repo *StatRepository) GenerateCSV(year, month int) (string, error) {
	var reportData []struct {
		ServiceId    int     `gorm:"column:service_id"`
		TotalRevenue float64 `gorm:"column:total_revenue"`
	}
	query := `
SELECT 
	service_id,
	SUM(amount) as total_revenue
FROM reservations 
WHERE status = 'confirmed'
AND EXTRACT(YEAR FROM created_at) = ?
AND EXTRACT(MONTH FROM created_at) = ?
GROUP BY service_id
ORDER BY total_revenue DESC
`
	if err := repo.Database.DB.Raw(query, year, month).Scan(&reportData).Error; err != nil {
		return "", fmt.Errorf("нет подтвержденных транзакций за %d-%02d", year, month)
	}

	reposrtsDirectory := "reports"
	if err := os.MkdirAll(reposrtsDirectory, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории для отчетов: %w", err)
	}
	filename := fmt.Sprintf("report-%d-%02d.csv", year, month)
	filePath := filepath.Join(reposrtsDirectory, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	if err := writer.Write([]string{"название услуги", "общая сумма за период"}); err != nil {
		return "", fmt.Errorf("ошибка записи заголовка: %w", err)
	}
	for _, item := range reportData {
		record := []string{
			fmt.Sprintf("Услуга %d", item.ServiceId),
			fmt.Sprintf("%.2f", item.TotalRevenue),
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("ошибка записи данных %w", err)
		}
	}
	return filePath, nil
}
