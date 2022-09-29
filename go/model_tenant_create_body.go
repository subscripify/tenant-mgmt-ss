/*
 * Private Services API  - OpenAPI 3.0
 *
 * This is the set of api endpoints to support access to Subscripify Super Services. Super Services are those that are available to Subscripify super tenants. They provide higher capabilities to super tenants to manage across multiple main tenants within the context of the Super-Tenant.<br><br> These APIs act as the front end to Private Services engineered and hosted by super tenants that need to access the subscripify platform.<br><br>Examples of super services available only to Super Tenants through this API are- <br>Tenant Management<br>Subscription and Plan Management <br>Billing and Payments<br><br> Subscripify also maintains a set of api endpoints to support access to Subscripify Public Services. Subscripify public services are services required by all tenants and users on the subscripify platform. All tenants, regardless of type have access to these services.<br><br>Examples of private services available to Super Tenants and to Main Tenants are- <br>Identity Services<br>Usage Analytics Services<br>For more information about Subscripify tenant architecture click here.
 *
 * API version: 0.0.1
 * Contact: william.ohara@subscripify.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package tenantapi

type TenantCreateBody struct {
	// Indicate which type of tenant to establish, main or super.
	TenantType string `json:"tenantType,omitempty"`
	// The alias name does not need to be unique and is used for quick reference when searching in UI. No starting spaces and no special characters.
	TenantAlias string `json:"tenantAlias,omitempty"`
	// The subdomain name string which used for the services namespace of the tenant and providing unique url for each tenant
	Subdomain string `json:"subdomain,omitempty"`
	// The services config UUID to use for a super tenant. Must be a valid services config UUID. This value must be empty when creating a main tenant.
	SuperServicesConfig string `json:"superServicesConfig,omitempty"`
	// The services config UUID to use for the tenant's public services. Must be a valid public services UUID.
	PublicServicesConfig string `json:"publicServicesConfig,omitempty"`
	// The private access config UUID to use for the tenant's public services. Must be a valid private access UUID. This value must be empty when creating a main tenant.
	PrivateAccessConfig string `json:"privateAccessConfig,omitempty"`
	// The public access config UUID to use for the tenant's public services. Must be a valid public access UUID.
	CustomAccessConfig string `json:"customAccessConfig,omitempty"`
}
