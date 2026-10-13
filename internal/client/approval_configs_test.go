// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestGetApprovalConfig(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/generic-approvals/v1/config/r1/ROLE" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("X-SailPoint-Experimental") != "true" {
			t.Errorf("missing experimental header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"r1","scope":"DOMAIN_OBJECT","escalationConfig":{"enabled":true,"escalationChain":[{"tier":1,"identityId":"","identityType":"MANAGER_OF"}]},"fallbackApprover":{"identityID":"g1","type":"GOVERNANCE_GROUP"}}`))
	})

	got, err := c.GetApprovalConfig(context.Background(), "r1", "ROLE")
	if err != nil {
		t.Fatalf("GetApprovalConfig: %v", err)
	}
	if got.ID == nil || *got.ID != "r1" || got.Scope == nil || *got.Scope != "DOMAIN_OBJECT" {
		t.Errorf("id/scope = %v/%v", got.ID, got.Scope)
	}
	if got.EscalationConfig == nil || got.EscalationConfig.EscalationChain == nil || len(*got.EscalationConfig.EscalationChain) != 1 {
		t.Fatalf("escalation chain = %+v", got.EscalationConfig)
	}
	if (*got.EscalationConfig.EscalationChain)[0].IdentityType != "MANAGER_OF" {
		t.Errorf("identity type = %q", (*got.EscalationConfig.EscalationChain)[0].IdentityType)
	}
	if got.FallbackApprover == nil || got.FallbackApprover.IdentityID != "g1" {
		t.Errorf("fallback = %+v", got.FallbackApprover)
	}
}

func TestPutApprovalConfigBody(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/generic-approvals/v1/config/ap1/ACCESS_PROFILE" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("X-SailPoint-Experimental") != "true" {
			t.Errorf("missing experimental header")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		for _, k := range []string{"id", "scope", "timeoutConfig"} {
			if _, ok := body[k]; ok {
				t.Errorf("body must not contain %q", k)
			}
		}
		fb, ok := body["fallbackApprover"].(map[string]any)
		if !ok {
			t.Fatalf("fallbackApprover missing: %v", body)
		}
		if _, ok := fb["identityID"]; !ok {
			t.Errorf("fallbackApprover has no identityID key: %v", fb)
		}
		rc, ok := body["reminderConfig"].(map[string]any)
		if !ok {
			t.Fatalf("reminderConfig missing: %v", body)
		}
		if v, ok := rc["enabled"]; !ok || v != false {
			t.Errorf("reminderConfig.enabled = %v, want present false", v)
		}
		sc, ok := body["serialChain"].([]any)
		if !ok || len(sc) != 0 {
			t.Errorf("serialChain = %v, want empty array", body["serialChain"])
		}
		writeJSON(t, w, http.StatusOK, map[string]any{})
	})

	disabled := false
	empty := []ApprovalChainItemAPI{}
	_, err := c.PutApprovalConfig(context.Background(), "ap1", "ACCESS_PROFILE", &ApprovalConfigAPI{
		ReminderConfig:   &ApprovalReminderConfigAPI{Enabled: &disabled},
		SerialChain:      &empty,
		FallbackApprover: &ApprovalFallbackApproverAPI{IdentityID: "g1", Type: "GOVERNANCE_GROUP"},
	})
	if err != nil {
		t.Fatalf("PutApprovalConfig: %v", err)
	}
}

func TestDeleteApprovalConfig(t *testing.T) {
	for _, tc := range []struct {
		status  int
		wantErr string
	}{
		{http.StatusNoContent, ""},
		{http.StatusNotFound, ""},
		{http.StatusForbidden, "access denied (403)"},
	} {
		c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete || r.URL.Path != "/generic-approvals/v1/config/r1/ROLE" {
				t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			}
			w.WriteHeader(tc.status)
		})
		err := c.DeleteApprovalConfig(context.Background(), "r1", "ROLE")
		switch {
		case tc.wantErr == "" && err != nil:
			t.Errorf("status %d: unexpected error %v", tc.status, err)
		case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
			t.Errorf("status %d: error = %v, want %q", tc.status, err, tc.wantErr)
		}
	}
}

func TestGetApprovalConfigNotFound(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := c.GetApprovalConfig(context.Background(), "r1", "ROLE")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestApprovalConfigEmptyArgs(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(_ http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected HTTP call %s %s", r.Method, r.URL.Path)
	})
	ctx := context.Background()
	if _, err := c.GetApprovalConfig(ctx, "", "ROLE"); err == nil {
		t.Error("empty id: expected error")
	}
	if _, err := c.GetApprovalConfig(ctx, "r1", ""); err == nil {
		t.Error("empty scope: expected error")
	}
}
