package main

import (
	"strings"
	"testing"
)

func TestBuildEntriesQuery(t *testing.T) {
	base := `audit_type:"login" tenant:"default"`
	got := buildEntriesQuery(base, 20, 40)
	want := `audit_type:"login" tenant:"default" | sort by (_time desc) | offset 40 | limit 20`
	if got != want {
		t.Fatalf("unexpected query:\nwant: %s\ngot:  %s", want, got)
	}
}

func TestNormalUserCannotOverrideAuditScope(t *testing.T) {
	query := buildLogsQL(TypeOperation, QueryParams{Tenant: "other", Username: "bob"}, UserContext{Tenant: "default", Username: "alice"})
	if !strings.Contains(query, `tenant:"default"`) || !strings.Contains(query, `username:"alice"`) || strings.Contains(query, `tenant:"other"`) || strings.Contains(query, `username:"bob"`) {
		t.Fatalf("audit scope escaped: %s", query)
	}
}
