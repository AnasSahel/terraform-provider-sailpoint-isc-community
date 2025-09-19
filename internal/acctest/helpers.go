// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// CheckListNotEmpty returns a function that verifies a list attribute is not empty
func CheckListNotEmpty(resourceName, listAttribute string) func(*terraform.State) error {
	return func(state *terraform.State) error {
		resource, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		// Check that list exists and has at least 1 item
		countAttr := listAttribute + ".#"
		countValue := resource.Primary.Attributes[countAttr]
		if countValue == "" {
			return fmt.Errorf("%s attribute not found", countAttr)
		}

		// Parse the count and validate it's a valid number > 0 (not empty)
		count, err := strconv.Atoi(countValue)
		if err != nil {
			return fmt.Errorf("expected %s to be a number, got %s", countAttr, countValue)
		}

		if count <= 0 {
			return fmt.Errorf("expected %s list to not be empty (%s > 0), but got %d", listAttribute, countAttr, count)
		}

		return nil
	}
}
