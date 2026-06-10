#!/bin/bash
# Import an existing source attribute sync config by source ID.
# Note: terraform destroy is a no-op for this resource — the config cannot be deleted via the API.
terraform import sailpoint_source_attribute_sync_config.ad "REPLACE_WITH_SOURCE_ID"
