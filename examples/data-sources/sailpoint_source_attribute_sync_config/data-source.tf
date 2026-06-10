data "sailpoint_source_attribute_sync_config" "ad" {
  source_id = sailpoint_source.ad.id
}

output "synced_attributes" {
  value = [
    for attr in data.sailpoint_source_attribute_sync_config.ad.attributes :
    attr.name if attr.enabled
  ]
}
