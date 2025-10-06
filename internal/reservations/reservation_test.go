package reservations_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task/internal/reservations"
	"task/internal/user"
	"task/pkg/dbhelper"
	"task/pkg/server"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestReservation_Duplicate(t *testing.T) {
	db := dbhelper.InitTestDb(t)
	defer dbhelper.CleanupTestDb(t, db)

	ts := httptest.NewServer(server.AppWithDB(db))
	defer ts.Close()

	testUser := user.CreateTestUser(t, ts, "TestUser")
	user.AddBalanceTestUser(ts, testUser.UserId, 200)

	reserveData, _ := json.Marshal(reservations.ReserveRequest{
		UserId:    testUser.UserId,
		OrderId:   1,
		ServiceId: 1,
		Amount:    100,
	})
	resp, _ := http.Post(ts.URL+"/reserve", "application/json", bytes.NewReader(reserveData))
	assert.Equal(t, 200, resp.StatusCode)

	resp, _ = http.Post(ts.URL+"/reserve", "application/json", bytes.NewReader(reserveData))
	assert.Equal(t, 400, resp.StatusCode)
}
