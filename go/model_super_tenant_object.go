/*
 * Private Services API  - OpenAPI 3.0
 *
 * This is the set of api endpoints to support access to Subscripify Super Services. Super Services are those that are available to Subscripify's super tenants. They provide higher capabilities to super tenants to manage across multiple main tenants within the context of the Super-Tenant.<br><br>  These APIs act as the front end to Private Services engineered and hosted by super tenants that need to access the subscripify platform.<br><br> Examples of super services available only to Super Tenants through this API are: <br> Tenant Management<br> Subscription and Plan Management <br> Billing and Payments<br><br> Subscripify also maintains a set of api endpoints to support access to Subscripify Public Services. Subscripify public services are services required by all tenants and users on the subscripify platform. All tenants, regardless of type have access to these services.<br><br> Examples of private services available to Super Tenants and to Main Tenants are: <br> Identity Services<br> Usage Analytics Services<br>  For more information about Subscripify's tenant architecture click here.
 *
 * API version: 0.0.1
 * Contact: william.ohara@subscripify.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package tenantapi

type SuperTenantObject struct {

	TenantUUID string `json:"tenantUUID,omitempty"`
	// The name of the organization at time of tenant creation and the alias used for searching by org name.
	OrgName string `json:"orgName,omitempty"`
	// The name prefix for the Kubernetes namespaces and cloud resources that make up this tenant.
	KubeNamespace string `json:"kubeNamespace,omitempty"`
	// The UUID of this tenant's owner tenant
	LiegeUUID string `json:"liegeUUID,omitempty"`

	CreateTimestamp string `json:"createTimestamp,omitempty"`
	// The services config UUID to use for the super-tenant. Must be a valid services config UUID
	SuperServicesConfig string `json:"superServicesConfig"`
	// Indicate which type of tenant one is spinning up. Main Tenants require only a public services config.  Super tenants require a super services config and a public config. Lord Tenants require super, public and internal services  configs.
	TenantType string `json:"tenantType,omitempty"`
	// The subdomain name string which used for the services namespace of the tenant and  providing unique url for each tenant
	Subdomain string `json:"subdomain"`
	// The services config UUID to use for the main-tenant's public services. The services config UUID  used must be a publicServices UUID and belong to the liege tenant
	PublicServicesConfig string `json:"publicServicesConfig"`
	// The cloud provider to deploy to. e.g. The only cloud provider supported (currently) is azure.
	SubscripifyDeploymentCloudLocation string `json:"subscripifyDeploymentCloudLocation"`
}