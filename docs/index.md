# NHN Cloud Terraform Provider

This document describes how to use NHN Cloud with Terraform.

## Example Usage

Before using Terraform, create a provider configuration file as follows.

The name of the provider file can be set randomly. This example uses `provider.tf` as the filename.

```
# Define required providers
terraform {
required_version = ">= 1.0.0"
  required_providers {
    nhncloud = {
      source  = "terraform.local/local/nhncloud"
      version = "1.0.0"
    }
  }
}

# Configure the nhncloud Provider
provider "nhncloud" {
  user_name   = "terraform-guide@nhncloud.com"
  tenant_id   = "aaa4c0a12fd84edeb68965d320d17129"
  password    = "difficultpassword"
  auth_url    = "https://api-identity-infrastructure.nhncloudservice.com/v2.0"
  region      = "KR1"
}
```

## Argument Reference

* `user_name` - (Required) Use the NHN Cloud ID.
* `tenant_id` - (Required) From **Compute > Instance > Management** on NHN Cloud console, click **Set API Endpoint** to check the Tenant ID.
* `password` - (Required) Use **API Password** that you saved in **Set API Endpoint**. Regarding how to set API passwords, see **User Guide > Compute > Instance > API Preparations**.
* `auth_url` - (Required) Specify the address of the NHN Cloud identification service. From **Compute > Instance > Management** on NHN Cloud console, click **Set API Endpoint** to check Identity URL.
* `region` - (Required) Enter the region to manage NHN Cloud resources.
    * `KR1`: Korea (Pangyo) Region
    * `KR2`: Korea (Pyeongchon) Region
    * `JP1`: Japan (Tokyo) Region

On the path where the provider configuration file is located, use the `init` command to initialize Terraform.

```sh
$ terraform init
```