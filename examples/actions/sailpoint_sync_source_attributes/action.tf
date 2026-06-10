# Manual trigger: run terraform apply and select this action to trigger a sync.
# Requires Terraform >= 1.14.
action "sailpoint_sync_source_attributes" "ad_resync" {
  config {
    source_id = sailpoint_source.ad.id
  }
}
