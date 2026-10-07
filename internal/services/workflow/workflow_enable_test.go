// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

type fakeWorkflowUpdater struct {
	calls   int
	gotID   string
	gotBody *client.WorkflowAPI
	err     error
}

func (f *fakeWorkflowUpdater) UpdateWorkflow(_ context.Context, id string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error) {
	f.calls++
	f.gotID = id
	f.gotBody = workflow
	if f.err != nil {
		return nil, f.err
	}
	out := *workflow
	return &out, nil
}

func boolPtr(v bool) *bool { return &v }

func TestEnableCreatedWorkflow_EnablesDisabledWorkflow(t *testing.T) {
	t.Parallel()

	updater := &fakeWorkflowUpdater{}
	created := &client.WorkflowAPI{ID: "wf-1", Name: "Repro", Enabled: boolPtr(false)}

	got, err := enableCreatedWorkflow(context.Background(), updater, created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updater.calls != 1 || updater.gotID != "wf-1" {
		t.Fatalf("expected one update for wf-1, got %d calls for %q", updater.calls, updater.gotID)
	}
	if updater.gotBody.Enabled == nil || !*updater.gotBody.Enabled {
		t.Error("expected update body to carry enabled = true")
	}
	if updater.gotBody.Name != "Repro" {
		t.Errorf("expected update body to keep the created workflow's fields, got name %q", updater.gotBody.Name)
	}
	if created.Enabled == nil || *created.Enabled {
		t.Error("expected the created workflow passed in to stay unmodified")
	}
	if got.Enabled == nil || !*got.Enabled {
		t.Error("expected returned workflow to be enabled")
	}
}

func TestEnableCreatedWorkflow_NilEnabledIsTreatedAsDisabled(t *testing.T) {
	t.Parallel()

	updater := &fakeWorkflowUpdater{}
	if _, err := enableCreatedWorkflow(context.Background(), updater, &client.WorkflowAPI{ID: "wf-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updater.calls != 1 {
		t.Errorf("expected one update, got %d", updater.calls)
	}
}

func TestEnableCreatedWorkflow_AlreadyEnabledSkipsUpdate(t *testing.T) {
	t.Parallel()

	updater := &fakeWorkflowUpdater{}
	created := &client.WorkflowAPI{ID: "wf-1", Enabled: boolPtr(true)}

	got, err := enableCreatedWorkflow(context.Background(), updater, created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updater.calls != 0 {
		t.Errorf("expected no update, got %d", updater.calls)
	}
	if got != created {
		t.Error("expected the created workflow to be returned as-is")
	}
}

func TestEnableCreatedWorkflow_PropagatesUpdateError(t *testing.T) {
	t.Parallel()

	updater := &fakeWorkflowUpdater{err: errors.New("boom")}
	_, err := enableCreatedWorkflow(context.Background(), updater, &client.WorkflowAPI{ID: "wf-1", Enabled: boolPtr(false)})
	if err == nil || err.Error() != "boom" {
		t.Errorf("expected update error to propagate, got %v", err)
	}
}
