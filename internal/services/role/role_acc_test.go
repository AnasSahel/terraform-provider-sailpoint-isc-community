// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_test

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
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
)

// providerFactories wires the in-process provider into the acceptance test
// harness.
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// fakeRoleAPI is an in-memory stand-in for the /roles/v1 endpoints. Like ISC,
// it fills in defaults on create (segments, dimensionRefs, accessRequestConfig)
// and a top-level replace overwrites the whole key, so a replace of
// accessRequestConfig without dimensionSchema drops the schema.
type fakeRoleAPI struct {
	t      *testing.T
	mu     sync.Mutex
	roles  map[string]map[string]any
	nextID int
}

func newFakeRoleServer(t *testing.T) (*httptest.Server, *fakeRoleAPI) {
	t.Helper()
	f := &fakeRoleAPI{t: t, roles: map[string]map[string]any{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	return srv, f
}

func (f *fakeRoleAPI) role(id string) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.roles[id]
}

func (f *fakeRoleAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/oauth/token" {
		writeJSON(w, map[string]any{"access_token": "fake", "token_type": "bearer", "expires_in": 3600})
		return
	}
	// /roles/v1[/{id}]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "roles" || parts[1] != "v1" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 2 {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var role map[string]any
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.nextID++
		id := fmt.Sprintf("role-%d", f.nextID)
		role["id"] = id
		for k, v := range map[string]any{
			"segments": []any{}, "dimensionRefs": []any{}, "enabled": false, "requestable": false,
			"dimensional": false, "created": "2026-01-01T00:00:00Z", "modified": "2026-01-01T00:00:00Z",
		} {
			if _, ok := role[k]; !ok {
				role[k] = v
			}
		}
		if role["dimensional"] == true {
			// ISC fills dimensionRefs from the role's dimensions.
			role["dimensionRefs"] = []any{map[string]any{"type": "DIMENSION", "id": "dimension-1", "name": "Dimension One"}}
		}
		if owner, ok := role["owner"].(map[string]any); ok {
			owner["name"] = "Owner"
		}
		if arc, ok := role["accessRequestConfig"].(map[string]any); ok {
			for k, v := range map[string]any{"reauthorizationRequired": false, "requireEndDate": false, "approvalSchemes": []any{}} {
				if _, set := arc[k]; !set {
					arc[k] = v
				}
			}
		} else {
			role["accessRequestConfig"] = map[string]any{
				"commentsRequired": nil, "denialCommentsRequired": nil, "approvalSchemes": []any{},
				"reauthorizationRequired": false, "requireEndDate": false,
				"maxPermittedAccessDuration": nil, "formDefinitionId": nil,
			}
		}
		f.roles[id] = role
		writeJSON(w, role)
		return
	}

	id := parts[2]
	role, ok := f.roles[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, role)
	case http.MethodPatch:
		var ops []client.JSONPatchOperation
		if err := json.NewDecoder(r.Body).Decode(&ops); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, op := range ops {
			key := strings.TrimPrefix(op.Path, "/")
			if key == "dimensional" || key == "dimensionRefs" {
				f.t.Errorf("PATCH %s: the roles API does not patch it", op.Path)
				http.Error(w, "not patchable: "+op.Path, http.StatusBadRequest)
				return
			}
			switch op.Op {
			case "replace", "add":
				role[key] = op.Value
			case "remove":
				delete(role, key)
			}
		}
		role["modified"] = "2026-01-02T00:00:00Z"
		writeJSON(w, role)
	case http.MethodDelete:
		delete(f.roles, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func roleConfig(baseURL, dynamicExtra, standardExtra string) string {
	return fmt.Sprintf(`
provider "sailpoint" {
  base_url      = %q
  client_id     = "fake"
  client_secret = "fake"
}

resource "sailpoint_role" "dynamic" {
  name        = "Dynamic Role"
  dimensional = true
  owner = {
    type = "IDENTITY"
    id   = "identity-1"
  }
%s
}

resource "sailpoint_role" "standard" {
  name = "Standard Role"
  owner = {
    type = "IDENTITY"
    id   = "identity-1"
  }
%s
}
`, baseURL, dynamicExtra, standardExtra)
}

// checkStoredDimensionSchema verifies the fake API still holds the dynamic
// role's dimension schema, i.e. no PATCH stripped it.
func checkStoredDimensionSchema(f *fakeRoleAPI) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["sailpoint_role.dynamic"]
		if !ok {
			return fmt.Errorf("sailpoint_role.dynamic not found in state")
		}
		arc, ok := f.role(rs.Primary.ID)["accessRequestConfig"].(map[string]any)
		if !ok {
			return fmt.Errorf("role %s has no accessRequestConfig", rs.Primary.ID)
		}
		if _, ok := arc["dimensionSchema"].(map[string]any); !ok {
			return fmt.Errorf("role %s lost its dimensionSchema: %v", rs.Primary.ID, arc)
		}
		return nil
	}
}

// TestAccRole_fakeAPI_dimensionSchema locks in #212: creating roles without
// segments or access_request_config is consistent, a dimension schema and form
// survive changes to other access request settings, an unconfigured
// dimensional does not force a replace, and changing dimensional does.
//
// It runs against an in-process fake API, so it needs TF_ACC=1 and a
// Terraform CLI but no SailPoint tenant.
func TestAccRole_fakeAPI_dimensionSchema(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}
	srv, f := newFakeRoleServer(t)
	const dynamic = "sailpoint_role.dynamic"
	const standard = "sailpoint_role.standard"
	const withSchema = `
  access_request_config = {
    dimension_schema = {
      dimension_attributes = [
        { name = "department", display_name = "Department" },
        { name = "location", derived = false },
      ]
    }
    form_definition_id = "form-def-1"
  }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// No segments, no access_request_config: consistent after apply, empty second plan,
				// also for a role whose dimension_refs the API fills.
				Config: roleConfig(srv.URL, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dynamic, "segments.#", "0"),
					resource.TestCheckResourceAttr(dynamic, "dimension_refs.#", "1"),
					resource.TestCheckResourceAttr(standard, "access_request_config.reauthorization_required", "false"),
				),
			},
			{
				Config: roleConfig(srv.URL, withSchema, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(dynamic, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dynamic, "access_request_config.dimension_schema.dimension_attributes.#", "2"),
					resource.TestCheckResourceAttr(dynamic, "access_request_config.dimension_schema.dimension_attributes.0.derived", "true"),
					resource.TestCheckResourceAttr(dynamic, "access_request_config.dimension_schema.dimension_attributes.1.derived", "false"),
					resource.TestCheckResourceAttr(dynamic, "access_request_config.form_definition_id", "form-def-1"),
				),
			},
			{
				// Schema and form left out: they are kept on the role.
				Config: roleConfig(srv.URL, `
  access_request_config = {
    comments_required = true
  }`, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(dynamic, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dynamic, "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr(dynamic, "access_request_config.dimension_schema.dimension_attributes.#", "2"),
					resource.TestCheckResourceAttr(dynamic, "access_request_config.form_definition_id", "form-def-1"),
					checkStoredDimensionSchema(f),
				),
			},
			{
				// The whole block left out: everything is kept, nothing to plan.
				Config:   roleConfig(srv.URL, "", ""),
				PlanOnly: true,
			},
			{
				// Unconfigured dimensional does not force a replace.
				Config: roleConfig(srv.URL, `
  access_request_config = {
    comments_required = true
  }`, `
  description = "Changed"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(standard, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(dynamic, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: strings.Replace(roleConfig(srv.URL, `
  access_request_config = {
    comments_required = true
  }`, `
  description = "Changed"`), "dimensional = true", "dimensional = false", 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(dynamic, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}
