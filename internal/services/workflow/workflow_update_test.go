// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

// fakeWorkflowServer mirrors one workflow on the server. Like ISC, it rejects
// a content update while the workflow is enabled.
type fakeWorkflowServer struct {
	current     client.WorkflowAPI
	puts        []client.WorkflowAPI // body of each UpdateWorkflow call
	updateErrOn int                  // if > 0, fail this UpdateWorkflow call (1-based)
}

func (f *fakeWorkflowServer) GetWorkflow(_ context.Context, _ string) (*client.WorkflowAPI, error) {
	out := f.current
	return &out, nil
}

func (f *fakeWorkflowServer) UpdateWorkflow(_ context.Context, _ string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error) {
	f.puts = append(f.puts, *workflow)
	if f.updateErrOn == len(f.puts) {
		return nil, errors.New("boom")
	}
	if isEnabled(f.current.Enabled) && workflow.Name != f.current.Name {
		return nil, errors.New("cannot update an enabled workflow")
	}
	f.current = *workflow
	out := f.current
	return &out, nil
}

func isEnabled(b *bool) bool { return b != nil && *b }

func TestUpdateWorkflow_DisabledWorkflowSinglePut(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{current: client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(false)}}
	got, err, restoreErr := updateWorkflow(context.Background(), srv, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(true)})
	if err != nil || restoreErr != nil {
		t.Fatalf("unexpected errors: %v, %v", err, restoreErr)
	}
	if len(srv.puts) != 1 {
		t.Fatalf("expected 1 PUT, got %d", len(srv.puts))
	}
	if got.Name != "new" || !isEnabled(got.Enabled) {
		t.Errorf("expected enabled workflow named new, got %+v", got)
	}
}

func TestUpdateWorkflow_EnabledWorkflowDisablesUpdatesReEnables(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{current: client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(true)}}
	got, err, restoreErr := updateWorkflow(context.Background(), srv, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(true)})
	if err != nil || restoreErr != nil {
		t.Fatalf("unexpected errors: %v, %v", err, restoreErr)
	}
	if len(srv.puts) != 3 {
		t.Fatalf("expected 3 PUTs (disable, update, enable), got %d", len(srv.puts))
	}
	if srv.puts[0].Name != "old" || isEnabled(srv.puts[0].Enabled) {
		t.Errorf("first PUT should disable the unchanged workflow, got %+v", srv.puts[0])
	}
	if srv.puts[1].Name != "new" || isEnabled(srv.puts[1].Enabled) {
		t.Errorf("second PUT should apply the update while disabled, got %+v", srv.puts[1])
	}
	if srv.puts[2].Name != "new" || !isEnabled(srv.puts[2].Enabled) {
		t.Errorf("third PUT should re-enable the updated workflow, got %+v", srv.puts[2])
	}
	if got.Name != "new" || !isEnabled(got.Enabled) {
		t.Errorf("expected enabled workflow named new, got %+v", got)
	}
}

func TestUpdateWorkflow_EnabledToDisabledSkipsReEnable(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{current: client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(true)}}
	got, err, restoreErr := updateWorkflow(context.Background(), srv, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(false)})
	if err != nil || restoreErr != nil {
		t.Fatalf("unexpected errors: %v, %v", err, restoreErr)
	}
	if len(srv.puts) != 2 {
		t.Fatalf("expected 2 PUTs (disable, update), got %d", len(srv.puts))
	}
	if got.Name != "new" || isEnabled(got.Enabled) {
		t.Errorf("expected disabled workflow named new, got %+v", got)
	}
}

func TestUpdateWorkflow_FailedUpdateReEnablesOriginal(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{
		current:     client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(true)},
		updateErrOn: 2,
	}
	got, err, restoreErr := updateWorkflow(context.Background(), srv, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(true)})
	if err == nil {
		t.Fatal("expected the update error")
	}
	if restoreErr != nil {
		t.Fatalf("unexpected restore error: %v", restoreErr)
	}
	if got != nil {
		t.Errorf("expected no workflow on failed update, got %+v", got)
	}
	if srv.current.Name != "old" || !isEnabled(srv.current.Enabled) {
		t.Errorf("expected the original workflow back and enabled, got %+v", srv.current)
	}
}

func TestUpdateWorkflow_FailedRestoreIsReported(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{current: client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(true)}}
	// Fail the update and the restore: every PUT after the disable fails.
	wrapped := &failFrom{fakeWorkflowServer: srv, from: 2}
	_, err, restoreErr := updateWorkflow(context.Background(), wrapped, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(true)})
	if err == nil || restoreErr == nil {
		t.Fatalf("expected both errors, got %v, %v", err, restoreErr)
	}
}

func TestUpdateWorkflow_FailedReEnableReturnsUpdatedWorkflow(t *testing.T) {
	t.Parallel()

	srv := &fakeWorkflowServer{
		current:     client.WorkflowAPI{ID: "wf-1", Name: "old", Enabled: boolPtr(true)},
		updateErrOn: 3,
	}
	got, err, restoreErr := updateWorkflow(context.Background(), srv, "wf-1", &client.WorkflowAPI{Name: "new", Enabled: boolPtr(true)})
	if err == nil {
		t.Fatal("expected the re-enable error")
	}
	if restoreErr != nil {
		t.Fatalf("unexpected restore error: %v", restoreErr)
	}
	if got == nil || got.Name != "new" || isEnabled(got.Enabled) {
		t.Errorf("expected the updated, disabled workflow so state matches the server, got %+v", got)
	}
}

// failFrom fails every UpdateWorkflow call from call number `from` onwards.
type failFrom struct {
	*fakeWorkflowServer
	from int
}

func (f *failFrom) UpdateWorkflow(ctx context.Context, id string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error) {
	if len(f.puts)+1 >= f.from {
		f.puts = append(f.puts, *workflow)
		return nil, errors.New("boom")
	}
	return f.fakeWorkflowServer.UpdateWorkflow(ctx, id, workflow)
}
