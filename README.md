Terraform Provider for NHN Cloud
============================

Requirements
------------

This is an example in `macOS / Apple silicon` architecture.

* [Terraform](https://www.terraform.io/downloads.html) 1.0.x
```sh
$ terraform version
Terraform v1.0.0
on darwin_arm64
```

* [Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)
```sh
$ go version
go version go1.20.5 darwin/arm64
```


Building The Provider
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

Please see the documentation at [registry.terraform.io]().

Or you can also check how to use Terraform in the NHN Cloud user guide [here](https://docs.nhncloud.com/ko/Compute/Instance/ko/terraform-guide/).


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.20+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
$ make build
```
