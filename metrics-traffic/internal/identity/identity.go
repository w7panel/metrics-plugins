package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type User struct {
	Username    string `json:"username"`
	Tenant      string `json:"tenant"`
	UserMode    string `json:"userMode"`
	IsAdmin     bool   `json:"isAdmin"`
	LogsEnabled bool   `json:"logsEnabled"`
}

var client = &http.Client{Timeout: 5 * time.Second}

// Verify asks PanelAuth for identity; request headers and query parameters cannot assert a role.
func Verify(r *http.Request) (User, error) {
	token := strings.TrimSpace(r.Header.Get("X-W7Panel-Token"))
	if token == "" {
		token = strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	}
	if token == "" {
		return User{}, errors.New("missing panel token")
	}
	base := strings.TrimRight(os.Getenv("PANEL_API_URL"), "/")
	if base == "" {
		base = "http://w7panel-offline.default.svc.cluster.local:8000"
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, base+"/panel-api/v1/metrics/plugin-identity", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("X-W7Panel-Token", token)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("panel identity rejected: HTTP %d", resp.StatusCode)
	}
	var payload struct {
		Data User `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<16)).Decode(&payload); err != nil {
		return User{}, err
	}
	if payload.Data.Username == "" || payload.Data.Tenant == "" {
		return User{}, errors.New("incomplete panel identity")
	}
	return payload.Data, nil
}
