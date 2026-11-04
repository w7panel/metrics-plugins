package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/w7panel/metrics-plugins/internal/identity"
)

func TestNormalUserCannotOverrideNamespace(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/traffic/apps?namespace=other", nil)
	params, err := parseTrafficParams(r, identity.User{Username: "alice", Tenant: "default", UserMode: "normal"})
	if err != nil || params.Namespace != "k3k-alice" {
		t.Fatalf("namespace = %q, %v", params.Namespace, err)
	}
}
