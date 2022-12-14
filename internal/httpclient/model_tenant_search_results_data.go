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

type TenantSearchResultsData struct {

	TenantUUID string `json:"tenantUUID,omitempty"`
	// The name of the organization at time of tenant creation and the alias used for searching by org name.
	TenantAlias string `json:"tenantAlias,omitempty"`
	// the subdomain of the tenant
	TenantSubdomain string `json:"tenantSubdomain,omitempty"`
	// the secondary domain of the tenant - does not include the top level domain
	TenantSecDomain string `json:"tenantSecDomain,omitempty"`
	// the top level domain of the tenant
	TenantTld string `json:"tenantTld,omitempty"`
	// Indicate which type of tenant. Lord tenants will see both \"super\" and \"main\" tenant types. Super tenants will only see \"main\" tenant types.
	TenantType string `json:"tenantType,omitempty"`
	// the lord config alias of the tenant - non lord tenants will be blank
	LordConfigAlias string `json:"lordConfigAlias,omitempty"`
	// the lord config UUID of the tenant - non lord tenants will be blank
	LordConfigUUID string `json:"lordConfigUUID,omitempty"`
	// the super config alias of the tenant - non super tenants will be blank
	SuperConfigAlias string `json:"superConfigAlias,omitempty"`
	// the super config UUID of the tenant - non super tenants will be blank
	SuperConfigUUID string `json:"superConfigUUID,omitempty"`
	// the public config alias of the tenant
	PublicConfigAlias string `json:"publicConfigAlias,omitempty"`
	// the public config UUID of the tenant
	PublicConfigUUID string `json:"publicConfigUUID,omitempty"`
	// the private access config (if applicable) alias of the tenant
	PrivateAccessAlias string `json:"privateAccessAlias,omitempty"`
	// the private access config (if applicable) UUID of the tenant
	PrivateAccessUUID string `json:"privateAccessUUID,omitempty"`
	// the custom access config (if applicable) alias of the tenant
	CustomAccessAlias string `json:"customAccessAlias,omitempty"`
	// the custom access config (if applicable) UUID of the tenant
	CustomAccessUUID string `json:"customAccessUUID,omitempty"`
}
