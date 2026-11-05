// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
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

// cvProviderFactories wires the in-process provider into the acceptance test
// harness.
var cvProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// fakeCriteriaValuesAPI is an in-memory stand-in for the /roles/v1 endpoints.
// Like ISC, every response carries stringValue and values on each criteria
// node, null when unset.
type fakeCriteriaValuesAPI struct {
	mu     sync.Mutex
	roles  map[string]map[string]any
	nextID int
}

func newFakeCriteriaValuesServer(t *testing.T) (*httptest.Server, *fakeCriteriaValuesAPI) {
	t.Helper()
	f := &fakeCriteriaValuesAPI{roles: map[string]map[string]any{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	return srv, f
}

func (f *fakeCriteriaValuesAPI) role(id string) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.roles[id]
}

func (f *fakeCriteriaValuesAPI) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.roles)
}

// seedUIRole stores a role shaped like one the SailPoint UI saved.
func (f *fakeCriteriaValuesAPI) seedUIRole(id string) {
	leaf := func(op, prop string, sv any, vals any) map[string]any {
		return map[string]any{
			"operation": op, "key": map[string]any{"type": "IDENTITY", "property": prop, "sourceId": nil},
			"stringValue": sv, "values": vals, "children": nil,
		}
	}
	role := map[string]any{
		"id": id, "name": "UI Role", "description": nil,
		"enabled": false, "requestable": false, "dimensional": false,
		"owner":    map[string]any{"type": "IDENTITY", "id": "identity-1", "name": "Owner"},
		"segments": []any{},
		"accessRequestConfig": map[string]any{
			"commentsRequired": false, "denialCommentsRequired": false,
			"reauthorizationRequired": false, "requireEndDate": false, "approvalSchemes": []any{},
		},
		"membership": map[string]any{"type": "STANDARD", "criteria": map[string]any{
			"operation": "AND", "key": nil, "stringValue": nil, "values": nil, "children": []any{
				leaf("EQUALS", "attribute.department", nil, []any{"ZZ-a", "ZZ-b"}),
				leaf("NOT_EQUALS", "attribute.location", nil, []any{"ZZ-c"}),
				map[string]any{"operation": "OR", "key": nil, "stringValue": nil, "values": nil, "children": []any{
					leaf("EQUALS", "attribute.x", nil, []any{"ZZ-d"}),
					leaf("EQUALS", "attribute.y", "ZZ-e", nil),
				}},
			},
		}},
		"created": "2026-01-01T00:00:00Z", "modified": "2026-01-01T00:00:00Z",
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.roles[id] = role
}

// cvNormaliseCriteria adds null stringValue/values keys where they are missing,
// as ISC does.
func cvNormaliseCriteria(role map[string]any) {
	m, ok := role["membership"].(map[string]any)
	if !ok {
		return
	}
	var walk func(node any)
	walk = func(node any) {
		n, ok := node.(map[string]any)
		if !ok {
			return
		}
		for _, k := range []string{"stringValue", "values"} {
			if _, set := n[k]; !set {
				n[k] = nil
			}
		}
		if children, ok := n["children"].([]any); ok {
			for _, c := range children {
				walk(c)
			}
		}
	}
	walk(m["criteria"])
}

func (f *fakeCriteriaValuesAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/oauth/token" {
		cvWriteJSON(w, map[string]any{"access_token": "fake", "token_type": "bearer", "expires_in": 3600})
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
		role["created"] = "2026-01-01T00:00:00Z"
		role["modified"] = "2026-01-01T00:00:00Z"
		// ISC always returns these; an empty segments list is dropped by omitempty on the way in.
		for k, v := range map[string]any{"segments": []any{}, "enabled": false, "requestable": false, "dimensional": false} {
			if _, ok := role[k]; !ok {
				role[k] = v
			}
		}
		if owner, ok := role["owner"].(map[string]any); ok {
			owner["name"] = "Owner"
		}
		cvNormaliseCriteria(role)
		f.roles[id] = role
		cvWriteJSON(w, role)
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
		cvWriteJSON(w, role)
	case http.MethodPatch:
		var ops []client.JSONPatchOperation
		if err := json.NewDecoder(r.Body).Decode(&ops); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, op := range ops {
			key := strings.TrimPrefix(op.Path, "/")
			switch op.Op {
			case "replace", "add":
				role[key] = op.Value
			case "remove":
				delete(role, key)
			}
		}
		cvNormaliseCriteria(role)
		role["modified"] = "2026-01-02T00:00:00Z"
		cvWriteJSON(w, role)
	case http.MethodDelete:
		delete(f.roles, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func cvWriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// cvRoleConfig renders the provider, the created role with the given two
// level-2 leaf bodies, and extra configuration.
func cvRoleConfig(baseURL, leaf1, leaf2, extra string) string {
	return fmt.Sprintf(`
provider "sailpoint" {
  base_url      = %q
  client_id     = "fake"
  client_secret = "fake"
}

locals {
  arc = {
    comments_required        = false
    denial_comments_required = false
    reauthorization_required = false
    require_end_date         = false
  }
}

resource "sailpoint_role" "created" {
  name                  = "Criteria Values Role"
  enabled               = false
  requestable           = false
  segments              = []
  access_request_config = local.arc
  owner = {
    type = "IDENTITY"
    id   = "identity-1"
  }
  membership = {
    type = "STANDARD"
    criteria = {
      operation = "AND"
      children = [
        {
          operation = "EQUALS"
          key       = { type = "IDENTITY", property = "attribute.department" }
          %s
        },
        {
          operation = "NOT_EQUALS"
          key       = { type = "IDENTITY", property = "attribute.location" }
          %s
        },
      ]
    }
  }
}
%s
`, baseURL, leaf1, leaf2, extra)
}

// cvUIRoleResource matches the role seeded by seedUIRole.
const cvUIRoleResource = `
resource "sailpoint_role" "ui" {
  name                  = "UI Role"
  enabled               = false
  requestable           = false
  segments              = []
  access_request_config = local.arc
  owner = {
    type = "IDENTITY"
    id   = "identity-1"
  }
  membership = {
    type = "STANDARD"
    criteria = {
      operation = "AND"
      children = [
        {
          operation = "EQUALS"
          key       = { type = "IDENTITY", property = "attribute.department" }
          values    = ["ZZ-a", "ZZ-b"]
        },
        {
          operation = "NOT_EQUALS"
          key       = { type = "IDENTITY", property = "attribute.location" }
          values    = ["ZZ-c"]
        },
        {
          operation = "OR"
          children = [
            {
              operation = "EQUALS"
              key       = { type = "IDENTITY", property = "attribute.x" }
              values    = ["ZZ-d"]
            },
            {
              operation    = "EQUALS"
              key          = { type = "IDENTITY", property = "attribute.y" }
              string_value = "ZZ-e"
            },
          ]
        },
      ]
    }
  }
}
`

// cvCheckStoredValues verifies the fake API holds the given values on each
// level-2 leaf of the created role, and no stringValue on values leaves.
func cvCheckStoredValues(f *fakeCriteriaValuesAPI, want ...[]any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["sailpoint_role.created"]
		if !ok {
			return fmt.Errorf("sailpoint_role.created not found in state")
		}
		raw, err := json.Marshal(f.role(rs.Primary.ID)["membership"])
		if err != nil {
			return err
		}
		var m struct {
			Criteria struct {
				Children []map[string]any `json:"children"`
			} `json:"criteria"`
		}
		if err := json.Unmarshal(raw, &m); err != nil {
			return err
		}
		if len(m.Criteria.Children) != len(want) {
			return fmt.Errorf("stored criteria has %d children, want %d: %s", len(m.Criteria.Children), len(want), raw)
		}
		for i, w := range want {
			leaf := m.Criteria.Children[i]
			if w == nil {
				continue
			}
			if !reflect.DeepEqual(leaf["values"], w) {
				return fmt.Errorf("stored leaf %d values = %v, want %v", i, leaf["values"], w)
			}
			if leaf["stringValue"] != nil {
				return fmt.Errorf("stored leaf %d has stringValue %v", i, leaf["stringValue"])
			}
		}
		return nil
	}
}

// TestAccRole_fakeAPI_criteriaValues locks in #212 item 3: criteria values
// lists create, import from a UI-built role with an empty plan, keep their
// order on update, mix with string_value, and are validated.
//
// It runs against an in-process fake API, so it needs TF_ACC=1 and a
// Terraform CLI but no SailPoint tenant.
func TestAccRole_fakeAPI_criteriaValues(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}
	srv, f := newFakeCriteriaValuesServer(t)
	f.seedUIRole("role-ui")

	const created = "sailpoint_role.created"
	const ab = `values = ["ZZ-a", "ZZ-b"]`
	const c = `values = ["ZZ-c"]`
	withUI := cvRoleConfig(srv.URL, ab, c, cvUIRoleResource)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: cvProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			if n := f.count(); n != 0 {
				return fmt.Errorf("%d roles left in the fake API", n)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: cvRoleConfig(srv.URL, ab, c, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(created, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(created, "membership.criteria.children.0.values.#", "2"),
					resource.TestCheckResourceAttr(created, "membership.criteria.children.0.values.1", "ZZ-b"),
					resource.TestCheckNoResourceAttr(created, "membership.criteria.children.0.string_value"),
					cvCheckStoredValues(f, []any{"ZZ-a", "ZZ-b"}, []any{"ZZ-c"}),
				),
			},
			{
				// Import a role the UI built: every leaf but one uses values.
				Config:             withUI,
				ResourceName:       "sailpoint_role.ui",
				ImportState:        true,
				ImportStateId:      "role-ui",
				ImportStatePersist: true,
			},
			{
				// After import, a matching configuration plans nothing.
				Config:   withUI,
				PlanOnly: true,
			},
			{
				// Reorder and change values: one in-place update, order kept.
				Config: cvRoleConfig(srv.URL, `values = ["ZZ-c2", "ZZ-a"]`, `values = ["ZZ-d", "ZZ-c"]`,
					cvUIRoleResource+`
data "sailpoint_role" "ui" {
  id = "role-ui"
}
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(created, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("sailpoint_role.ui", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(created, "membership.criteria.children.0.values.0", "ZZ-c2"),
					resource.TestCheckResourceAttr(created, "membership.criteria.children.0.values.1", "ZZ-a"),
					resource.TestCheckResourceAttr(created, "membership.criteria.children.1.values.0", "ZZ-d"),
					resource.TestCheckResourceAttr(created, "membership.criteria.children.1.values.1", "ZZ-c"),
					cvCheckStoredValues(f, []any{"ZZ-c2", "ZZ-a"}, []any{"ZZ-d", "ZZ-c"}),
					resource.TestCheckResourceAttr("data.sailpoint_role.ui", "membership.criteria.children.0.values.#", "2"),
					resource.TestCheckResourceAttr("data.sailpoint_role.ui", "membership.criteria.children.0.values.1", "ZZ-b"),
					resource.TestCheckResourceAttr("data.sailpoint_role.ui", "membership.criteria.children.2.children.1.string_value", "ZZ-e"),
				),
			},
			{
				// Both on one leaf is rejected at plan time. It is not the last
				// step, so the post-test destroy runs with a valid configuration.
				Config:      cvRoleConfig(srv.URL, `string_value = "ZZ-s"`+"\n          "+ab, c, cvUIRoleResource),
				ExpectError: regexp.MustCompile("not both"),
			},
			{
				// Mixed: one leaf string_value, one values.
				Config: cvRoleConfig(srv.URL, `string_value = "ZZ-s"`, `values = ["ZZ-d", "ZZ-c"]`, cvUIRoleResource),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(created, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(created, "membership.criteria.children.0.string_value", "ZZ-s"),
					resource.TestCheckNoResourceAttr(created, "membership.criteria.children.0.values.#"),
					cvCheckStoredValues(f, nil, []any{"ZZ-d", "ZZ-c"}),
				),
			},
		},
	})
}
