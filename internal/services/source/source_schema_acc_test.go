// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

// fakeSchemaAPI is an in-memory stand-in for the /v2025/sources/{id}/schemas
// endpoints. It returns attributes in the reverse of the order they were sent,
// the way ISC may reorder a schema, so tests prove the provider does not
// depend on attribute order.
type fakeSchemaAPI struct {
	mu      sync.Mutex
	schemas map[string]client.SourceSchemaAPI
	nextID  int
}

func newFakeSchemaServer(t *testing.T) *httptest.Server {
	t.Helper()
	f := &fakeSchemaAPI{schemas: map[string]client.SourceSchemaAPI{}}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)
	return srv
}

func (f *fakeSchemaAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/oauth/token" {
		writeJSON(w, map[string]any{"access_token": "fake", "token_type": "bearer", "expires_in": 3600})
		return
	}

	// /v2025/sources/{sourceId}/schemas[/{schemaId}]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 || parts[0] != "v2025" || parts[1] != "sources" || parts[3] != "schemas" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 4 {
		switch r.Method {
		case http.MethodGet:
			out := []client.SourceSchemaAPI{}
			names := r.URL.Query().Get("include-names")
			for _, s := range f.schemas {
				if names == "" || s.Name == names {
					out = append(out, reversed(s))
				}
			}
			writeJSON(w, out)
		case http.MethodPost:
			var s client.SourceSchemaAPI
			if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			f.nextID++
			s.ID = fmt.Sprintf("schema-%d", f.nextID)
			f.schemas[s.ID] = s
			writeJSON(w, reversed(s))
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id := parts[4]
	s, ok := f.schemas[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, reversed(s))
	case http.MethodPut:
		var upd client.SourceSchemaAPI
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		upd.ID = id
		f.schemas[id] = upd
		writeJSON(w, reversed(upd))
	case http.MethodDelete:
		delete(f.schemas, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func reversed(s client.SourceSchemaAPI) client.SourceSchemaAPI {
	s.Attributes = slices.Clone(s.Attributes)
	slices.Reverse(s.Attributes)
	return s
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func sourceSchemaConfig(baseURL, attributes string) string {
	return fmt.Sprintf(`
provider "sailpoint" {
  base_url      = %q
  client_id     = "fake"
  client_secret = "fake"
}

resource "sailpoint_source_schema" "test" {
  source_id          = "source-1"
  name               = "account"
  native_object_type = "User"
  identity_attribute = "sAMAccountName"

  attributes = [
%s
  ]
}
`, baseURL, attributes)
}

// TestAccSourceSchema_attributesOrderInsensitive locks in #165 and #167: once
// applied, a config that lists the same attributes in another order and omits
// the computed flags plans nothing, and adding one attribute is an in-place
// update that keeps the other attributes' flags.
//
// It runs against an in-process fake API, so it needs TF_ACC=1 and a
// Terraform CLI but no SailPoint tenant.
func TestAccSourceSchema_attributesOrderInsensitive(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests (uses a fake API, needs the Terraform CLI)")
	}

	srv := newFakeSchemaServer(t)
	const resourceName = "sailpoint_source_schema.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: sourceSchemaConfig(srv.URL, `
    { name = "sAMAccountName", type = "STRING", description = "Login" },
    { name = "mail", type = "STRING", description = "Email" },
    { name = "memberOf", type = "STRING", description = "Groups", is_multi = true, is_entitlement = true, is_group = true },`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "attributes.*", map[string]string{
						"name": "memberOf", "is_multi": "true", "is_entitlement": "true", "is_group": "true",
					}),
				),
			},
			{
				// Same set, different order, flags omitted: no diff.
				Config: sourceSchemaConfig(srv.URL, `
    { name = "memberOf", type = "STRING", description = "Groups" },
    { name = "mail", type = "STRING", description = "Email" },
    { name = "sAMAccountName", type = "STRING", description = "Login" },`),
				PlanOnly: true,
			},
			{
				// One added attribute: an in-place update, existing flags kept.
				Config: sourceSchemaConfig(srv.URL, `
    { name = "memberOf", type = "STRING", description = "Groups" },
    { name = "mail", type = "STRING", description = "Email" },
    { name = "sAMAccountName", type = "STRING", description = "Login" },
    { name = "department", type = "STRING", description = "Department" },`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "attributes.*", map[string]string{
						"name": "memberOf", "is_multi": "true", "is_entitlement": "true", "is_group": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "attributes.*", map[string]string{
						"name": "department", "is_multi": "false",
					}),
				),
			},
		},
	})
}
