// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

// providerFactories wires the in-process provider into the acceptance test
// harness.
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

const fakeRoleID = "role-1"

// patchablePaths are the fields the dimensions API accepts in a JSON Patch.
var patchablePaths = map[string]bool{
	"/name": true, "/description": true, "/owner": true,
	"/accessProfiles": true, "/entitlements": true, "/membership": true,
}

// fakeDimensionAPI is an in-memory stand-in for /roles/v1/{roleId}/dimensions.
// Like ISC, it keeps the owner as sent with a null name, fills access profile
// and entitlement names, stores [] for missing lists, writes "" on AND and OR
// nodes, and returns access profiles in a different order on every response.
type fakeDimensionAPI struct {
	t              *testing.T
	mu             sync.Mutex
	dims           map[string]map[string]any
	nextID         int
	responses      int
	lastPatchPaths []string
}

func newFakeDimensionServer(t *testing.T) (*httptest.Server, *fakeDimensionAPI) {
	t.Helper()
	f := &fakeDimensionAPI{t: t, dims: map[string]map[string]any{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	return srv, f
}

func (f *fakeDimensionAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/oauth/token" {
		writeJSON(w, map[string]any{"access_token": "fake", "token_type": "bearer", "expires_in": 3600})
		return
	}

	// /roles/v1/role-1/dimensions[/{id}]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 || len(parts) > 5 || parts[0] != "roles" || parts[1] != "v1" || parts[2] != fakeRoleID || parts[3] != "dimensions" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 4 {
		switch r.Method {
		case http.MethodPost:
			var dim map[string]any
			if err := json.NewDecoder(r.Body).Decode(&dim); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			f.nextID++
			id := fmt.Sprintf("dimension-%d", f.nextID)
			dim["id"] = id
			dim["parentId"] = fakeRoleID
			dim["created"] = "2026-01-01T00:00:00Z"
			dim["modified"] = "2026-01-01T00:00:00Z"
			for _, k := range []string{"owner", "accessProfiles", "entitlements", "membership"} {
				dim[k] = normalise(k, dim[k])
			}
			f.dims[id] = dim
			writeJSON(w, f.response(dim))
		case http.MethodGet:
			list := make([]any, 0, len(f.dims))
			for _, dim := range f.dims {
				list = append(list, f.response(dim))
			}
			writeJSON(w, list)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id := parts[4]
	dim, ok := f.dims[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, f.response(dim))
	case http.MethodPatch:
		if ct := r.Header.Get("Content-Type"); ct != "application/json-patch+json" {
			f.t.Errorf("PATCH Content-Type = %q, want application/json-patch+json", ct)
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}
		var ops []client.JSONPatchOperation
		if err := json.NewDecoder(r.Body).Decode(&ops); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.lastPatchPaths = nil
		for _, op := range ops {
			if !patchablePaths[op.Path] {
				f.t.Errorf("PATCH %s: the dimensions API does not patch it", op.Path)
				http.Error(w, "not patchable: "+op.Path, http.StatusBadRequest)
				return
			}
			f.lastPatchPaths = append(f.lastPatchPaths, op.Path)
			key := strings.TrimPrefix(op.Path, "/")
			switch op.Op {
			case "replace", "add":
				dim[key] = normalise(key, op.Value)
			case "remove":
				dim[key] = normalise(key, nil)
			default:
				f.t.Errorf("PATCH %s: unexpected op %q", op.Path, op.Op)
			}
		}
		dim["modified"] = "2026-01-02T00:00:00Z"
		writeJSON(w, f.response(dim))
	case http.MethodDelete:
		delete(f.dims, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// normalise stores a field the way ISC does.
func normalise(key string, v any) any {
	switch key {
	case "owner":
		if owner, ok := v.(map[string]any); ok {
			owner["name"] = nil
			return owner
		}
		return nil
	case "accessProfiles", "entitlements":
		list, _ := v.([]any)
		out := make([]any, 0, len(list))
		for _, item := range list {
			if ref, ok := item.(map[string]any); ok {
				ref["name"] = fmt.Sprintf("Name %v", ref["id"])
				out = append(out, ref)
			}
		}
		return out
	case "membership":
		if m, ok := v.(map[string]any); ok {
			if c, ok := m["criteria"].(map[string]any); ok {
				emptyConnectorValues(c)
			}
			return m
		}
		return nil
	}
	return v
}

// emptyConnectorValues writes "" on AND and OR nodes, as the UI does.
func emptyConnectorValues(node map[string]any) {
	if op := node["operation"]; op == "AND" || op == "OR" {
		node["stringValue"] = ""
	}
	children, _ := node["children"].([]any)
	for _, c := range children {
		if child, ok := c.(map[string]any); ok {
			emptyConnectorValues(child)
		}
	}
}

// response copies a stored dimension and rotates its access profiles by one
// more position on every call, so consecutive reads return a different order.
func (f *fakeDimensionAPI) response(dim map[string]any) map[string]any {
	out := make(map[string]any, len(dim))
	for k, v := range dim {
		out[k] = v
	}
	if aps, ok := dim["accessProfiles"].([]any); ok && len(aps) > 1 {
		f.responses++
		n := f.responses % len(aps)
		rotated := append(append([]any{}, aps[n:]...), aps[:n]...)
		out["accessProfiles"] = rotated
	}
	return out
}

func (f *fakeDimensionAPI) onlyDimension() (string, map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for id, dim := range f.dims {
		return id, dim
	}
	return "", nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func dimensionConfig(baseURL, name string, accessProfileIDs ...string) string {
	var aps strings.Builder
	for _, id := range accessProfileIDs {
		fmt.Fprintf(&aps, "    { type = \"ACCESS_PROFILE\", id = %q },\n", id)
	}
	return fmt.Sprintf(`
provider "sailpoint" {
  base_url      = %q
  client_id     = "fake"
  client_secret = "fake"
}

resource "sailpoint_role_dimension" "test" {
  role_id     = %q
  name        = %q
  description = "Managed by Terraform"
  owner = {
    type = "IDENTITY"
    id   = "identity-1"
  }
  access_profiles = [
%s  ]
  membership = {
    type = "STANDARD"
    criteria = {
      operation = "AND"
      children = [
        {
          operation    = "EQUALS"
          key          = { type = "IDENTITY", property = "attribute.department" }
          string_value = "Engineering"
        },
        {
          operation = "OR"
          children = [
            {
              operation    = "EQUALS"
              key          = { type = "IDENTITY", property = "attribute.location" }
              string_value = "London"
            },
            {
              operation    = "EQUALS"
              key          = { type = "IDENTITY", property = "attribute.location" }
              string_value = "Paris"
            },
          ]
        },
      ]
    }
  }
}

data "sailpoint_role_dimension" "test" {
  role_id = sailpoint_role_dimension.test.role_id
  id      = sailpoint_role_dimension.test.id
}
`, baseURL, fakeRoleID, name, aps.String())
}

func checkLastPatchPaths(f *fakeDimensionAPI, want ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if !reflect.DeepEqual(f.lastPatchPaths, want) {
			return fmt.Errorf("last PATCH paths = %v, want %v", f.lastPatchPaths, want)
		}
		return nil
	}
}

// TestAccRoleDimension_fakeAPI locks in #212 item 1: access profiles returned
// in a different order plan nothing, the empty "" ISC writes on AND/OR nodes
// does not drift, only changed fields are patched, import works, and a
// membership cleared in the UI is restored.
//
// It runs against an in-process fake API, so it needs TF_ACC=1 and a
// Terraform CLI but no SailPoint tenant.
func TestAccRoleDimension_fakeAPI(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}

	srv, f := newFakeDimensionServer(t)
	const res = "sailpoint_role_dimension.test"
	const ds = "data.sailpoint_role_dimension.test"
	full := dimensionConfig(srv.URL, "Engineering - London", "ap-1", "ap-2", "ap-3")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		CheckDestroy: func(_ *terraform.State) error {
			if id, _ := f.onlyDimension(); id != "" {
				return fmt.Errorf("dimension %s still exists", id)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// Create. The harness also checks the plan after apply is empty.
				Config: full,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "role_id", fakeRoleID),
					resource.TestCheckResourceAttr(res, "access_profiles.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(res, "access_profiles.*", map[string]string{"id": "ap-2", "name": "Name ap-2"}),
					resource.TestCheckResourceAttr(res, "owner.name", ""),
					resource.TestCheckNoResourceAttr(res, "membership.criteria.string_value"),
					resource.TestCheckNoResourceAttr(res, "membership.criteria.children.1.string_value"),
					resource.TestCheckResourceAttr(ds, "name", "Engineering - London"),
					resource.TestCheckResourceAttr(ds, "access_profiles.#", "3"),
				),
			},
			{
				// Every refresh sees a new access profile order: nothing to plan.
				Config:   full,
				PlanOnly: true,
			},
			{
				Config:   full,
				PlanOnly: true,
			},
			{
				// Only the changed fields are patched.
				Config: dimensionConfig(srv.URL, "Engineering - Paris", "ap-1", "ap-3"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "name", "Engineering - Paris"),
					resource.TestCheckResourceAttr(res, "access_profiles.#", "2"),
					checkLastPatchPaths(f, "/name", "/accessProfiles"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[res]
					if !ok {
						return "", fmt.Errorf("%s not found in state", res)
					}
					return fakeRoleID + "/" + rs.Primary.ID, nil
				},
			},
			{
				// A UI save clears the membership and changes the description: the next apply restores both.
				PreConfig: func() {
					_, dim := f.onlyDimension()
					f.mu.Lock()
					defer f.mu.Unlock()
					dim["membership"] = nil
					dim["description"] = "Changed in the UI"
				},
				Config: dimensionConfig(srv.URL, "Engineering - Paris", "ap-1", "ap-3"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(res, "membership.criteria.operation", "AND"),
					checkLastPatchPaths(f, "/description", "/membership"),
				),
			},
		},
	})
}
