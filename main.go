package main

import (
	"flag"

	"github.com/nhn-cloud/terraform-provider-nhncloud/nhncloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        debugMode,
		ProviderFunc: nhncloud.Provider,
	})
}
