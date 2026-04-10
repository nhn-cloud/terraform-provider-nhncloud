package nhncloud

import (
	"fmt"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/gophercloud/utils/terraform/auth"
)

// Use openstackbase.Config as the base/foundation of this provider's
// Config struct.
type Config struct {
	auth.Config

	// tenantClients caches ProviderClients keyed by tenant ID for cross-tenant operations.
	tenantClients      map[string]*gophercloud.ProviderClient
	tenantClientsMutex sync.Mutex
}

func (c *Config) NasStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(nasStorageClientInit, region, "nasv1")
}

// NasStorageV1ClientForTenant returns a NAS storage client scoped to the given tenant.
// If tenantID is empty or matches the provider's tenant, delegates to NasStorageV1Client.
func (c *Config) NasStorageV1ClientForTenant(region, tenantID string) (*gophercloud.ServiceClient, error) {
	if tenantID == "" || tenantID == c.TenantID {
		return c.NasStorageV1Client(region)
	}

	providerClient, err := c.getOrCreateTenantClient(tenantID)
	if err != nil {
		return nil, err
	}

	serviceClient, err := nasStorageClientInit(providerClient, gophercloud.EndpointOpts{
		Region:       c.DetermineRegion(region),
		Availability: clientconfig.GetEndpointType(c.EndpointType),
	})
	if err != nil {
		return nil, err
	}

	return c.DetermineEndpoint(serviceClient, "nasv1"), nil
}

func (c *Config) getOrCreateTenantClient(tenantID string) (*gophercloud.ProviderClient, error) {
	c.tenantClientsMutex.Lock()
	defer c.tenantClientsMutex.Unlock()

	if client, ok := c.tenantClients[tenantID]; ok {
		return client, nil
	}

	client, err := openstack.NewClient(c.IdentityEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating provider client for tenant %s: %w", tenantID, err)
	}

	client.HTTPClient = c.OsClient.HTTPClient
	client.UserAgent = c.OsClient.UserAgent
	client.Context = c.OsClient.Context
	client.MaxBackoffRetries = c.OsClient.MaxBackoffRetries
	client.RetryBackoffFunc = c.OsClient.RetryBackoffFunc

	authOpts := *c.AuthOpts
	authOpts.TenantID = tenantID
	authOpts.TenantName = ""

	if err := openstack.Authenticate(client, authOpts); err != nil {
		return nil, fmt.Errorf("error authenticating for tenant %s: %w", tenantID, err)
	}

	c.tenantClients[tenantID] = client
	return client, nil
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
