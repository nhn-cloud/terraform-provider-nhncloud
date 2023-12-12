Terraform Provider for NHN Cloud
============================

Requirements
------------

* [Terraform](https://www.terraform.io/downloads.html) 1.0.x

* [Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)

Building the Provider
---------------------

Clone the repository

```sh
$ git clone git@github.com:nhn/terraform-provider-nhncloud.git
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-nhncloud
$ make build
```

Provider Usage
-----------------

Please see the [NHN Cloud Terraform Provider documentation]() for how to use NHN Cloud Terraform Provider.

You can also check the [NHN Cloud Terraform User Guide](https://docs.nhncloud.com/en/Compute/Instance/en/terraform-guide/).


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.20+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
$ make build
```

