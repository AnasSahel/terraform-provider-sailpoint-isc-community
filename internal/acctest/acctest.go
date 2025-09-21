// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// This file is kept for future provider-level tests.
// Individual resource and data source tests are located in their respective service packages.

const (
	ProviderConfig = `
		provider "sailpoint" {}
	`
)

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func TestAccSailPointTransformsDataSource() string {
	return ProviderConfig + `
		data "sailpoint_transforms" "test" {}
	`
}

func TestAccSailPointTransformResourceCreate() string {
	return ProviderConfig + `
		resource "sailpoint_transform" "test" {
			name = "terraform-poc-001"
			type = "upper"
			attributes = jsonencode({
				"input" = "name"
			})
		}
	`
}

func TestAccSailPointTransformResourceUpdate() string {
	return ProviderConfig + `
		resource "sailpoint_transform" "test" {
			name = "terraform-poc-001"
			type = "upper"
			attributes = jsonencode({
				"input" = "updated name"
			})
		}
	`
}

func TestAccSailPointManagedClusterDataSource() string {
	return `
		terraform {
			required_providers {
				random = {
					source = "hashicorp/random"
				}
			}
		}
		provider "sailpoint" {}

		resource "sailpoint_managed_cluster" "dependency" {
			name = "tf-test-${random_id.cluster.hex}"
			type = "idn"
			description = "Dependency cluster for data source test"
			configuration = {
				gmt_offset = "-5"
			}
		}

		resource "random_id" "cluster" {
			byte_length = 4
		}

		data "sailpoint_managed_cluster" "test" {
			id = sailpoint_managed_cluster.dependency.id
		}
	`
}

func TestAccSailPointManagedClusterResourceCreate() string {
	return `
		terraform {
			required_providers {
				random = {
					source = "hashicorp/random"
				}
			}
		}
		provider "sailpoint" {}

		resource "random_id" "cluster" {
			byte_length = 4
		}

		resource "sailpoint_managed_cluster" "test" {
			name = "tf-test-${random_id.cluster.hex}"
			type = "idn"
			description = "Test managed cluster created by Terraform"
			configuration = {
				gmt_offset = "-5"
			}
		}
	`
}

func TestAccSailPointManagedClusterResourceUpdate() string {
	return `
		terraform {
			required_providers {
				random = {
					source = "hashicorp/random"
				}
			}
		}
		provider "sailpoint" {}

		resource "random_id" "cluster" {
			byte_length = 4
		}

		resource "sailpoint_managed_cluster" "test" {
			name = "tf-test-${random_id.cluster.hex}"
			type = "idn"
			description = "Updated test managed cluster description"
			configuration = {
				gmt_offset = "-5"
				debug_mode = "true"
			}
		}
	`
}

// Source resource test configurations
func TestAccSailPointSourcesDataSource() string {
	return ProviderConfig + `
		data "sailpoint_sources" "test" {}
	`
}

func TestAccSailPointSourceDataSource() string {
	return ProviderConfig + `
		resource "random_id" "source" {
			byte_length = 4
		}

		resource "sailpoint_source" "dependency" {
			name = "tf-test-dep-${random_id.source.hex}"
			description = "Dependency source for data source test"
			connector = "delimited-file"
			owner = jsonencode({
				"type" = "IDENTITY"
				"id" = "2c91808570313110017040b06f344ec9"
				"name" = "john.doe"
			})
			connector_attributes = jsonencode({
				"file" = "test.csv"
			})
		}

		data "sailpoint_source" "test" {
			id = sailpoint_source.dependency.id
		}
	`
}

func TestAccSailPointSourceResourceCreate() string {
	return ProviderConfig + `
		data "sailpoint_sources" "first" {}

		resource "sailpoint_source" "test" {
			name = "tf-test-delimited-source"
			description = "Test delimited source created by Terraform"
			connector = "delimited-file"
			owner = jsonencode({
				"type" = "IDENTITY"
				"id" = data.sailpoint_sources.first.sources[0].owner.id
				"name" = data.sailpoint_sources.first.sources[0].owner.name
			})
			connector_attributes = jsonencode({
				"file" = "users.csv"
			})
		}
	`
}

func TestAccSailPointSourceResourceUpdate() string {
	return ProviderConfig + `
		data "sailpoint_sources" "first" {}

		resource "sailpoint_source" "test" {
			name = "tf-test-delimited-source"
			description = "Updated test delimited source"
			connector = "delimited-file"
			owner = jsonencode({
				"type" = "IDENTITY"
				"id" = data.sailpoint_sources.first.sources[0].owner.id
				"name" = data.sailpoint_sources.first.sources[0].owner.name
			})
			connector_attributes = jsonencode({
				"file" = "users.csv"
				"delimiter" = ";"
			})
		}
	`
}

func TestAccSailPointConnectorsDataSource() string {
	return ProviderConfig + `
		data "sailpoint_connectors" "test" {}
	`
}
