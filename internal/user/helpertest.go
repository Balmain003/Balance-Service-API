package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetTestUserBalance(t *testing.T, ts *httptest.Server, userID int) (*User, error) {

	res, err := http.Get(fmt.Sprintf("%s/balance/%d", ts.URL, userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("expected 200 got %d. Response: %s", res.StatusCode, string(body))
	}

	var userData User
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &userData, nil
}

func AddBalanceTestUser(ts *httptest.Server, userID int, amount float64) (*User, error) {
	addBalanceData, err := json.Marshal(AddBalanceRequest{
		UserId: userID,
		Amount: amount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	req, err := http.NewRequest("PATCH", ts.URL+"/balance", bytes.NewReader(addBalanceData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("expected 200 got %d. Response: %s", res.StatusCode, string(body))
	}

	var response User
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &response, nil
}

func CreateTestUser(t *testing.T, ts *httptest.Server, name string) *User {
	userData, _ := json.Marshal(&UserCreateRequest{Name: name})
	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewReader(userData))
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdUser User
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &createdUser)
	resp.Body.Close()

	return &createdUser
}
