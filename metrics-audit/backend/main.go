package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/w7panel/metrics-plugins/metrics-audit/internal/identity"
)

func auditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := identity.Verify(r)
	if err != nil {
		identity.Error(w, http.StatusUnauthorized, err)
		return
	}
	if r.URL.Path == "/api/v1/audit/logs/status" {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		identity.JSON(w, http.StatusOK, CheckStatus(ctx, user.LogsEnabled))
		return
	}
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params := QueryParams{Page: page, PageSize: pageSize, Username: q.Get("username"), Tenant: q.Get("tenant"), Success: q.Get("success"), Method: q.Get("method"), Path: q.Get("path"), StartTime: q.Get("startTime"), EndTime: q.Get("endTime")}
	current := UserContext{Tenant: user.Tenant, Username: user.Username, UserMode: user.UserMode, IsAdmin: user.IsAdmin}
	var result QueryResult
	switch r.URL.Path {
	case "/api/v1/audit/login-logs":
		result, err = QueryLoginLogs(r.Context(), params, current)
	case "/api/v1/audit/operation-logs":
		result, err = QueryOperationLogs(r.Context(), params, current)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		identity.Error(w, http.StatusBadGateway, err)
		return
	}
	identity.JSON(w, http.StatusOK, result)
}

func main() {
	if os.Getenv("LOGS_BASE_URL") == "" {
		panic("LOGS_BASE_URL is required")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/audit/", auditHandler)
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
