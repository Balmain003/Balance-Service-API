package dbhelper

import (
	"os/user"
	"task/internal/history"
	"task/internal/reservations"
	"task/internal/stat"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitTestDb(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=192.168.0.103 user=postgres password=mypassword dbname=postgres_test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&user.User{}, &stat.ReportFile{}, &history.Transaction{}, &reservations.Reservations{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	t.Logf("Successfully migrated all tables")
	return db
}

func CleanupTestDb(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []interface{}{&user.User{}, &history.Transaction{}, &reservations.Reservations{}, &stat.ReportFile{}}

	for _, table := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			t.Logf("Warning: failed to clean table: %v", err)
		}
	}
}
