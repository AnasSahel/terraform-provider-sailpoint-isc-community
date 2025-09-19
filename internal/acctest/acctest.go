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
