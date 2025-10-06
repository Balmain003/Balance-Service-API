package user_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"task/internal/history"
	"task/internal/reservations"
	"task/internal/stat"
	"task/internal/user"
	"task/pkg/dbhelper"
	"task/pkg/server"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateSuccsess(t *testing.T) {
	db := dbhelper.InitTestDb(t)
	defer dbhelper.CleanupTestDb(t, db)
	ts := httptest.NewServer(server.AppWithDB(db))
	defer ts.Close()

	createdUser := user.CreateTestUser(t, ts, "TestUser")

	var dbUser user.User
	if err := db.Where("name = ?", "TestUser").First(&dbUser).Error; err != nil {
		t.Fatalf("User not found in database: %v", err)
	}
	if createdUser.Name != dbUser.Name {
		t.Fatalf("Expected name %s, got %s", createdUser.Name, dbUser.Name)
	}
}

func TestAddBalanceSuccess(t *testing.T) {
	db := dbhelper.InitTestDb(t)
	defer dbhelper.CleanupTestDb(t, db)
	if err := db.AutoMigrate(&user.User{}, &history.Transaction{}, &reservations.Reservations{}, &stat.ReportFile{}); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}
	ts := httptest.NewServer(server.AppWithDB(db))
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

	response, err := user.AddBalanceTestUser(ts, createdUser.UserId, 100)
	if err != nil {
		t.Fatalf("Failed to add balance: %v", err)
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

func TestTransferSuccsess(t *testing.T) {
	db := dbhelper.InitTestDb(t)
	defer dbhelper.CleanupTestDb(t, db)
	ts := httptest.NewServer(server.AppWithDB(db))
	defer ts.Close()

	user1 := user.CreateTestUser(t, ts, "User1")
	user2 := user.CreateTestUser(t, ts, "User2")

	user.AddBalanceTestUser(ts, user1.UserId, 200)

	transferData, _ := json.Marshal(user.TransferRequest{
		FromUserId: user1.UserId,
		ToUserId:   user2.UserId,
		Amount:     100,
	})
	resp, err := http.Post(ts.URL+"/transfer", "application/json", bytes.NewReader(transferData))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	user1Balance, _ := user.GetTestUserBalance(t, ts, user1.UserId)
	user2Balance, _ := user.GetTestUserBalance(t, ts, user2.UserId)

	assert.Equal(t, 100.0, user1Balance.Balance)
	assert.Equal(t, 100.0, user2Balance.Balance)
}
