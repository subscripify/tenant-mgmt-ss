/*
 * Private Services API  - OpenAPI 3.0
 *
 * This is the set of api endpoints to support access to Subscripify Super Services. Super Services are those that are available to Subscripify's super tenants. They provide higher capabilities to super tenants to manage across multiple main tenants within the context of the Super-Tenant.<br><br>  These APIs act as the front end to Private Services engineered and hosted by super tenants that need to access the subscripify platform.<br><br> Examples of super services available only to Super Tenants through this API are: <br> Tenant Management<br> Subscription and Plan Management <br> Billing and Payments<br><br> Subscripify also maintains a set of api endpoints to support access to Subscripify Public Services. Subscripify public services are services required by all tenants and users on the subscripify platform. All tenants, regardless of type have access to these services.<br><br> Examples of private services available to Super Tenants and to Main Tenants are: <br> Identity Services<br> Usage Analytics Services<br>  For more information about Subscripify's tenant architecture click here.
 *
 * API version: 0.0.1
 * Contact: william.ohara@subscripify.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"log"
	"net/http"

	tenantapi "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/go"
	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/subscripify-logger"
)

func main() {
	subscripifylogger.InfoLog.Printf("tenant management service started")

	router := tenantapi.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
