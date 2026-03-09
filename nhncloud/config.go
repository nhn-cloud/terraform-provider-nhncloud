package nhncloud

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/terraform/auth"
)

// Use openstackbase.Config as the base/foundation of this provider's
// Config struct.
type Config struct {
	auth.Config
}

func (c *Config) NasStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(nasStorageClientInit, region, "nasv1")
}

func nasStorageClientInit(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "nasv1")
}

func initClientOpts(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}
