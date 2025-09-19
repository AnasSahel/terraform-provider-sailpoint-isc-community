# Example of importing an existing SailPoint managed cluster
# 
# Step 1: Create the resource configuration
resource "sailpoint_managed_cluster" "example" {
  name        = "Existing Cluster Name"
  description = "Description of the existing cluster"
  type        = "idn"

  configuration = {
    # Add configuration to match the existing cluster
    gmt_offset = "-5"
  }

  lifecycle {
    # Prevent accidental deletion during import process
    prevent_destroy = true
  }
}

# Step 2: Run the import command
# terraform import sailpoint_managed_cluster.example "your-cluster-id-here"
#
# Step 3: Run terraform plan to see any differences
# terraform plan
#
# Step 4: Update the configuration above to match the imported state
# Step 5: Run terraform plan again to ensure no changes are detected
