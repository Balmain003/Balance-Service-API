package stat

import "time"

type Stat struct {
	ServiceId    int       `json:"sevice_id"`
	TotalRevenue float64   `json:"total_revenue"`
	Period       time.Time `json:"period"`
}

type ReportFile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName  string    `gorm:"not null" json:"file_name"`
	FilePath  string    `gorm:"not null" json:"file_path"`
	Period    time.Time `gorm:"not null" json:"period"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
