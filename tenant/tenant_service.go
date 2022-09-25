package tenant

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/tenantdb"
	"github.com/google/uuid"
)

type iHttpResponse interface {
	GetHttpResponse() *httpResponseData
	logAndGenerateHttpResponseData(responseCode int, message string, location string)
	logAndGenerateNewTenantResponse(nID uuid.UUID, location string)
}

type newTenantResponse struct {
	TenantUUID string
}

type httpResponseData struct {
	HttpResponseCode int
	Message          string
	NewTenant        newTenantResponse
}

func (r *httpResponseData) GetHttpResponse() *httpResponseData {
	return r
}

// this function generates the data for http responses, regardless if the api uses them in a response or not
func (hr *httpResponseData) logAndGenerateHttpResponseData(responseCode int, message string, location string) {

	hr.HttpResponseCode = responseCode
	hr.Message = message
	log.Printf(
		"tenant service generated http response : %v %s %s",
		responseCode,
		strings.ToLower(message),
		location,
	)

}

// this function generates a new tenant response and assigns it to a
func (tr *httpResponseData) logAndGenerateNewTenantResponse(nID uuid.UUID, location string) {

	tr.NewTenant.TenantUUID = nID.String()

	log.Printf(
		"tenant service generated tenant:%s %s",
		nID.String(),
		location,
	)

}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// this function creates a new lord tenant object and then attempts to insert the object into the database. if it fails it will generate an response
// structure interface for to use for generating an http response
func NewLordTenant(
	tenantAlias string,
	topLevelDomain string,
	secondaryDomain string,
	subdomain string,
	lordServicesConfig string,
	superServicesConfig string,
	publicServicesConfig string,
	cloudLocation CloudLocation,
	createdBy string) iHttpResponse {
	//no special processing required - this is a pass through to maintain the pattern. the NewTenant function covers the factory for non lord tenant types
	nlt, resp := createLordTenant(tenantAlias, topLevelDomain, secondaryDomain, subdomain, lordServicesConfig, superServicesConfig, publicServicesConfig, cloudLocation, createdBy)

	//if new lord tenant object is successful then try to insert into the database
	if !isNil(nlt) {
		//setting up a 10 second timeout (could be less)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tdb := tenantdbserv.Tdb.Handle

		insertStr := `INSERT INTO tenant (
		tenant_uuid, 
		tenant_alias,
		top_level_domain,
		secondary_domain,
		subdomain, 
		kube_namespace_prefix, 
		lord_services_config, 
		super_services_config,
		public_services_config,
		subscripify_deployment_cloud_location,
		is_lord_tenant,
		created_by
			)
		VALUES (UUID_TO_BIN(?), ?,?,?,?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?),UUID_TO_BIN(?),?,?,?);`

		_, rc, message := tenantdbserv.InsertResponseHelper(tdb.ExecContext(ctx, insertStr,
			nlt.getTenantUUID(),
			nlt.getAlias(),
			nlt.getTopLevelDomain(),
			nlt.getSecondaryDomainName(),
			nlt.getSubdomainName(),
			nlt.getKubeNamespacePrefix(),
			nlt.getLordServicesConfig(),
			nlt.getSuperServicesConfig(),
			nlt.getPublicServicesConfig(),
			nlt.getCloudLocation(),
			nlt.isLordTenant(),
			nlt.getTenantCreator()))

		//if the response is not a 200 then an insert error ocurred
		if rc != 200 {
			resp.logAndGenerateHttpResponseData(rc, message, "NewLordTenant")
		} else {
			resp.logAndGenerateHttpResponseData(200, "created new tenant", "NewLordTenant")
			resp.logAndGenerateNewTenantResponse(nlt.getTenantUUID(), "NewLordTenant")
		}
	}

	return resp
}

func NewTenant(
	tenantType string,
	tenantAlias string,
	subdomain string,
	superServicesConfig string,
	publicServicesConfig string,
	privateAccessConfig string,
	customAccessConfig string,
	liegeUUID string,
	createdBy string) iHttpResponse {
	// factory function for new super or main tenants. this function will also look up lord tenants using liege tenant value

	//setting up a 10 second timeout (could be less)

	if tenantType == string(MainTenant) {

		nmt, resp := createMainTenant(tenantAlias, subdomain, publicServicesConfig, customAccessConfig, liegeUUID, createdBy)

		if !isNil(nmt) {
			//setting up a 10 second timeout (could be less)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			tdb := tenantdbserv.Tdb.Handle

			insertStr := `INSERT INTO tenant (
			tenant_uuid, 
			tenant_alias,
			top_level_domain,
			secondary_domain,
			subscripify_deployment_cloud_location,
			subdomain, 
			kube_namespace_prefix, 
			public_services_config,
			custom_access_config,
		  created_by
				)
			VALUES (UUID_TO_BIN(?),
			 ?,
			 (SELECT top_level_domain, secondary_domain, subscripify_deployment_cloud_location FROM tenants WHERE tenant_uuid = ? LIMIT 1),  
			 ?, 
			 ?, 
			 UUID_TO_BIN(?), 
			 UUID_TO_BIN(?),
			 ?;`

			_, rc, message := tenantdbserv.InsertResponseHelper(tdb.ExecContext(ctx, insertStr,
				nmt.getTenantUUID(),
				nmt.getAlias(),
				nmt.getLiegeUUID(),
				nmt.getSubdomainName(),
				nmt.getKubeNamespacePrefix(),
				nmt.getPublicServicesConfig(),
				nmt.getCustomAccessConfig(),
				nmt.getTenantCreator()))

			//if the response is not a 200 then an insert error ocurred
			if rc != 200 {
				resp.logAndGenerateHttpResponseData(rc, message, "NewTenant:main")
			} else {
				resp.logAndGenerateHttpResponseData(200, "created new main tenant", "NewTenant:main")
				resp.logAndGenerateNewTenantResponse(nmt.getTenantUUID(), "NewTenant:main")
			}
		}

		return resp

	} else if tenantType == string(SuperTenant) {
		nst, resp := createSuperTenant(tenantAlias, subdomain, superServicesConfig, publicServicesConfig, privateAccessConfig, customAccessConfig, liegeUUID, createdBy)

		if !isNil(nst) {
			//setting up a 10 second timeout (could be less)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			tdb := tenantdbserv.Tdb.Handle

			insertStr := `INSERT INTO tenant (
			tenant_uuid, 
			tenant_alias,
			top_level_domain,
			secondary_domain,
			subscripify_deployment_cloud_location,
			subdomain, 
			kube_namespace_prefix,
			super_services_config, 
			public_services_config,
			private_access_config,
			custom_access_config,
		  created_by
				)
			VALUES (UUID_TO_BIN(?),
			 ?,
			 (SELECT top_level_domain, secondary_domain, subscripify_deployment_cloud_location FROM tenants WHERE tenant_uuid = ? LIMIT 1),  
			 ?, 
			 ?, 
			 UUID_TO_BIN(?),
			 UUID_TO_BIN(?),
			 UUID_TO_BIN(?), 
			 UUID_TO_BIN(?),
			 ?;`

			_, rc, message := tenantdbserv.InsertResponseHelper(tdb.ExecContext(ctx, insertStr,
				nst.getTenantUUID(),
				nst.getAlias(),
				nst.getLiegeUUID(),
				nst.getSubdomainName(),
				nst.getKubeNamespacePrefix(),
				nst.getSuperServicesConfig(),
				nst.getPublicServicesConfig(),
				nst.getPrivateAccessConfig(),
				nst.getCustomAccessConfig(),
				nst.getTenantCreator()))

			//if the response is not a 200 then an insert error ocurred
			if rc != 200 {
				resp.logAndGenerateHttpResponseData(rc, message, "NewTenant:main")
			} else {
				resp.logAndGenerateHttpResponseData(200, "created new main tenant", "NewTenant:main")
				resp.logAndGenerateNewTenantResponse(nst.getTenantUUID(), "NewTenant:main")
			}
		}

		return resp

	} else if tenantType == string(LordTenant) {
		var resp iHttpResponse
		resp.logAndGenerateHttpResponseData(405, "lord tenants need to be created using the NewLordTenant function a POST to /lord-tenants", "NewTenant:Lord")

		return resp
	}
	var resp iHttpResponse
	resp.logAndGenerateHttpResponseData(405, "this does not seem to be a valid tenant type - only super or main tenants can be created with this endpoint", "NewTenant:Lord")

	return resp
}
