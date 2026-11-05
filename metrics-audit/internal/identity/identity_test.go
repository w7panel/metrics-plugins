package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyUsesPanelIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/panel-api/v1/metrics/plugin-identity" || r.Header.Get("X-W7Panel-Token") != "valid" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(`{"data":{"username":"alice","tenant":"default","userMode":"normal","isAdmin":false}}`))
	}))
	defer server.Close()
	t.Setenv("PANEL_API_URL", server.URL)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/audit/login-logs?tenant=other", nil)
	request.Header.Set("X-W7Panel-Token", "valid")
	request.Header.Set("X-User-Role", "founder")
	user, err := Verify(request)
	if err != nil || user.Username != "alice" || user.IsAdmin {
		t.Fatalf("verified identity = %+v, %v", user, err)
	}
	request.Header.Set("X-W7Panel-Token", "invalid")
	if _, err := Verify(request); err == nil {
		t.Fatal("accepted invalid token")
	}
}
