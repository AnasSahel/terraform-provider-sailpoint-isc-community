# Schedule account aggregation for an existing source to run on weekday mornings.
# Requires the parent sailpoint_source to exist first.
resource "sailpoint_source_aggregation_schedule" "accounts_nightly" {
  source_id     = "REPLACE_WITH_SOURCE_ID"
  schedule_type = "ACCOUNTS"

  schedule = {
    type = "WEEKLY"

    hours = {
      type   = "LIST"
      values = ["2"] # 02:00 UTC
    }

    days = {
      type   = "LIST"
      values = ["MON", "TUE", "WED", "THU", "FRI"]
    }

    time_zone_id = "UTC"
  }

  enable_threshold = false
}

# Entitlement aggregation on Sunday nights with thresholding enabled.
# If more than 20% of entitlements change in a single run, the aggregation
# is paused and requires manual review in the SailPoint UI.
resource "sailpoint_source_aggregation_schedule" "entitlements_weekly" {
  source_id     = "REPLACE_WITH_SOURCE_ID"
  schedule_type = "ENTITLEMENTS"

  schedule = {
    type = "WEEKLY"

    hours = {
      type   = "LIST"
      values = ["1"] # 01:00 UTC
    }

    days = {
      type   = "LIST"
      values = ["SUN"]
    }

    time_zone_id = "UTC"
  }

  enable_threshold = true
  threshold        = 20
}

# Daily account aggregation — no days component needed for DAILY type.
resource "sailpoint_source_aggregation_schedule" "accounts_daily" {
  source_id     = "REPLACE_WITH_SOURCE_ID"
  schedule_type = "ACCOUNTS"

  schedule = {
    type = "DAILY"

    hours = {
      type   = "LIST"
      values = ["4"] # 04:00 UTC
    }

    time_zone_id = "America/Chicago"
  }

  enable_threshold = false
}

# Import an existing schedule that was configured outside Terraform:
#   terraform import sailpoint_source_aggregation_schedule.accounts_nightly REPLACE_WITH_SOURCE_ID/ACCOUNTS
