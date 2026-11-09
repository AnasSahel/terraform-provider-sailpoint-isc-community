// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package launcher_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// fakeLauncherAPI mimics /launchers/v1 the way ISC behaves: it stamps the
// owner with the calling identity on every write and returns refs without
// a name.
type fakeLauncherAPI struct {
	mu        sync.Mutex
	launchers map[string]client.LauncherAPI
	nextID    int
	writes    int
}

func newFakeLauncherServer(t *testing.T) *httptest.Server {
	t.Helper()
	f := &fakeLauncherAPI{launchers: map[string]client.LauncherAPI{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	return srv
}

func (f *fakeLauncherAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/oauth/token" {
		writeJSON(w, map[string]any{"access_token": "fake", "token_type": "bearer", "expires_in": 3600})
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "launchers" || parts[1] != "v1" {
		http.NotFound(w, r)
		return
	}

	decode := func() (client.LauncherAPI, bool) {
		var l client.LauncherAPI
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return l, false
		}
		// ISC overwrites the owner with the caller and does not echo names.
		l.Owner = &client.ObjectRefAPI{Type: "USER", ID: "caller-identity-id"}
		if l.Reference != nil {
			l.Reference.Name = ""
		}
		f.writes++
		l.Modified = fmt.Sprintf("2026-01-01T00:00:%02dZ", f.writes)
		return l, true
	}

	switch {
	case len(parts) == 2 && r.Method == http.MethodPost:
		l, ok := decode()
		if !ok {
			return
		}
		f.nextID++
		l.ID = fmt.Sprintf("launcher-%d", f.nextID)
		l.Created = "2026-01-01T00:00:00Z"
		f.launchers[l.ID] = l
		writeJSON(w, l)
	case len(parts) == 3:
		cur, exists := f.launchers[parts[2]]
		if !exists {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, cur)
		case http.MethodPut:
			l, ok := decode()
			if !ok {
				return
			}
			l.ID, l.Created = cur.ID, cur.Created
			f.launchers[l.ID] = l
			writeJSON(w, l)
		case http.MethodDelete:
			delete(f.launchers, parts[2])
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func launcherConfig(baseURL, ownerBlock string) string {
	return fmt.Sprintf(`
provider "sailpoint" {
  base_url      = %q
  client_id     = "fake"
  client_secret = "fake"
}

resource "sailpoint_launcher" "test" {
  name        = "Test launcher"
  description = "Launches a workflow"
  type        = "INTERACTIVE_PROCESS"
  config      = jsonencode({ workflowId = "wf-1" })
%s
  reference = {
    type = "WORKFLOW"
    id   = "wf-1"
  }
}
`, baseURL, ownerBlock)
}

const ownerBlock = `
  owner = {
    type = "USER"
    id   = "caller-identity-id"
  }
`

// TestAccLauncher_noPermanentDiffWithoutOwner reproduces issue #209: a
// launcher whose owner is left to the server (or removed from the config
// after an upgrade) must reach "No changes" on the next plan.
func TestAccLauncher_noPermanentDiffWithoutOwner(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}
	srv := newFakeLauncherServer(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Old-style config, with owner set (as in 3.1.x).
				Config: launcherConfig(srv.URL, ownerBlock),
			},
			{
				// Owner removed from config: nothing to change.
				Config: launcherConfig(srv.URL, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccLauncher_createWithoutOwner covers a fresh launcher with no owner.
func TestAccLauncher_createWithoutOwner(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}
	srv := newFakeLauncherServer(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{Config: launcherConfig(srv.URL, "")},
			{
				Config: launcherConfig(srv.URL, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}
