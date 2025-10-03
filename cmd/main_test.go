package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"task/internal/history"
	"task/internal/reservations"
	"task/internal/stat"
	"task/internal/user"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateSuccsess(t *testing.T) {
	db := InitTestDb(t)
	defer CleanupTestDb(t, db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, err := json.Marshal(&user.UserCreateRequest{
		Name: "TestUser",
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	res, err := http.Post(ts.URL+"/user", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal()
	}

	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Fatalf("Expected %d got %d", 201, res.StatusCode)
	}

	var response user.UserCreateResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	var dbUser user.User
	if err := db.Where("name = ?", "TestUser").First(&dbUser).Error; err != nil {
		t.Fatalf("User not found in database: %v", err)
	}
}

func TestAddBalanceSuccsess(t *testing.T) {
	db := InitTestDb(t)
	defer CleanupTestDb(t, db)
	if err := db.AutoMigrate(&user.User{}, &history.Transaction{}, &reservations.Reservations{}, &stat.ReportFile{}); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}
	ts := httptest.NewServer(App())
	defer ts.Close()

	createdUserData, err := json.Marshal(&user.UserCreateRequest{
		Name: "TestUser",
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	createdUserRes, err := http.Post(ts.URL+"/user", "application/json", bytes.NewReader(createdUserData))
	if err != nil {
		t.Fatal()
	}

	defer createdUserRes.Body.Close()

	if createdUserRes.StatusCode != 201 {
		t.Fatalf("Expected %d got %d", 201, createdUserRes.StatusCode)
	}

	var createdUser user.User

	createdBody, _ := io.ReadAll(createdUserRes.Body)
	if err := json.Unmarshal(createdBody, &createdUser); err != nil {
		t.Fatalf("Failed to unmarshal created user: %v", err)
	}

	addBalanseData, err := json.Marshal(user.AddBalanceRequest{
		UserId: createdUser.UserId,
		Amount: 100,
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	req, err := http.NewRequest("PATCH", ts.URL+"/balance", bytes.NewReader(addBalanseData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("Expected 200 got %d. Response: %s", res.StatusCode, string(body))
	}

	var response user.User
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Balance != 100 {
		t.Fatalf("Expected balance %f got %f", 100.0, response.Balance)
	}

	var dbUser user.User
	if err := db.Where("user_id = ?", createdUser.UserId).First(&dbUser).Error; err != nil {
		t.Fatalf("User not found in database: %v", err)
	}
	if dbUser.Balance != 100 {
		t.Fatalf("Expected balance in DB %f got %f", 100.0, dbUser.Balance)
	}
}

func InitTestDb(t *testing.T) *gorm.DB {
	t.Helper()

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
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
