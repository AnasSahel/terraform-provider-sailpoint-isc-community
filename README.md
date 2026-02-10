# Terraform Provider for SailPoint Identity Security Cloud (Community)

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-623CE4?style=flat&logo=terraform)](https://www.terraform.io)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

A Terraform provider for managing [SailPoint Identity Security Cloud (ISC)](https://www.sailpoint.com/) resources. This community-maintained provider enables infrastructure-as-code management of SailPoint ISC configurations.

**Current Version:** v2.0.0

## Implemented Resources

| Resource | Data Source | Description |
|----------|-------------|-------------|
| `sailpoint_identity_attribute` | `sailpoint_identity_attribute` | Manage identity attribute configurations |
| `sailpoint_transform` | `sailpoint_transform` | Manage SailPoint transforms for attribute manipulation |
| `sailpoint_form_definition` | `sailpoint_form_definition` | Manage form definitions for access requests and workflows |
| `sailpoint_workflow` | `sailpoint_workflow` | Manage workflow definitions and steps |
| `sailpoint_workflow_trigger` | - | Manage workflow triggers (EVENT, SCHEDULED, EXTERNAL) |
| `sailpoint_launcher` | `sailpoint_launcher` | Manage launchers to trigger workflows through the SailPoint UI |
| `sailpoint_lifecycle_state` | `sailpoint_lifecycle_state` | Manage lifecycle states within identity profiles |
| `sailpoint_source_schema` | `sailpoint_source_schema` | Manage source schema definitions for accounts and entitlements |
| `sailpoint_identity_profile` | `sailpoint_identity_profile` | Manage identity profiles and attribute mappings |

## SailPoint v2025 API Coverage

**Summary:** 9 of 83 API endpoints implemented (10.8%)

### ✅ Implemented APIs

| API Endpoint | Resource/Data Source |
|--------------|---------------------|
| Custom Forms | `sailpoint_form_definition` |
| Identity Attributes | `sailpoint_identity_attribute` |
| Identity Profiles | `sailpoint_identity_profile` |
| Launchers | `sailpoint_launcher` |
| Lifecycle States | `sailpoint_lifecycle_state` |
| Source Schemas | `sailpoint_source_schema` |
| Transforms | `sailpoint_transform` |
| Triggers | `sailpoint_workflow_trigger` |
| Workflows | `sailpoint_workflow` |

### ❌ Not Yet Implemented APIs

| API Endpoint | API Endpoint | API Endpoint |
|--------------|--------------|--------------|
| Access Model Metadata | Access Profiles | Access Request Approvals |
| Access Request Identity Metrics | Access Requests | Account Activities |
| Account Aggregations | Account Usages | Accounts |
| Api Usage | Application Discovery | Branding |
| Certification Campaign Filters | Certification Campaigns | Certification Summaries |
| Certifications | Classify Source | Configuration Hub |
| Connector Customizers | Connector Rule Management | Connectors |
| Custom Password Instructions | Custom User Levels | Data Access Security |
| Data Segmentation | Declassify Source | Dimensions |
| Entitlements | Global Tenant Security Settings | Governance Groups |
| IAI Access Request Recommendations | IAI Common Access | Identity History |
| Lifecycle States | Machine Account Classify | Machine Account Mappings |
| Machine Accounts | Machine Classification Config | Machine Identities |
| Managed Clients | Managed Cluster Types | Managed Clusters |
| MFA Configuration | Multi-Host Integration | Non-Employee Lifecycle Management |
| Notifications | OAuth Clients | Org Config |
| Parameter Storage | Password Configuration | Password Dictionary |
| Password Management | Requestable Objects | Role Insights |
| Roles | Saved Search | Scheduled Search |
| Search | Search Attribute Configuration | Segments |
| Service Desk Integration | SIM Integrations | SOD Policies |
| SOD Violations | Source Usages | Sources |
| SP-Config | Suggested Entitlement Description | Tagged Objects |
| Tags | Task Management | Tenant |
| Tenant Context | UI Metadata | Work Items |
| Work Reassignment | | |
