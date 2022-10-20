package tenant

import (
	"context"
	"fmt"
	"strings"
	"time"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/tenant-mgmt-ss.git/internal/tenantdb"
)

// func isNil(i interface{}) bool {
// 	return i == nil || reflect.ValueOf(i).IsNil()
// }

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
	var resp httpResponseData
	nlt, responseCode, err := createLordTenant(tenantAlias, topLevelDomain, secondaryDomain, subdomain, lordServicesConfig, superServicesConfig, publicServicesConfig, cloudLocation, createdBy)

	if err != nil {

		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())

	} else {
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

		_, rc, err := tenantdbserv.HttpResponseHelperSQLInsert(tdb.ExecContext(ctx, insertStr,
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
		if err != nil {
			resp.generateHttpResponseCodeAndMessage(rc, err.Error())

		} else {
			resp.generateHttpResponseCodeAndMessage(200, "created new tenant")
			resp.generateNewTenantResponse(nlt.getTenantUUID())
		}
	}

	return &resp
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

	var resp httpResponseData
	if tenantType == string(MainTenant) {

		nmt, responseCode, err := createMainTenant(tenantAlias, subdomain, publicServicesConfig, customAccessConfig, liegeUUID, createdBy)
		if err != nil {
			resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
			return &resp
		}
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
				liege_uuid,
				lord_uuid,
				created_by
					)   SELECT UUID_TO_BIN(?), ?, top_level_domain, secondary_domain, subscripify_deployment_cloud_location, ?, ?, UUID_TO_BIN(?),
				 UUID_TO_BIN(?), UUID_TO_BIN(?), lord_uuid, ? FROM tenant WHERE tenant_uuid = UUID_to_bin(?);`

		_, rc, err := tenantdbserv.HttpResponseHelperSQLInsert(tdb.ExecContext(ctx, insertStr,
			nmt.getTenantUUID(),
			nmt.getAlias(),
			nmt.getSubdomainName(),
			nmt.getKubeNamespacePrefix(),
			nmt.getPublicServicesConfig(),
			nmt.getCustomAccessConfig(),
			nmt.getLiegeUUID(),
			nmt.getTenantCreator(),
			nmt.getLiegeUUID()))

		//if the response is not a 200 then an insert error ocurred
		if err != nil {
			resp.generateHttpResponseCodeAndMessage(rc, err.Error())
			return &resp
		}
		resp.generateHttpResponseCodeAndMessage(200, "created new main tenant")
		resp.generateNewTenantResponse(nmt.getTenantUUID())

		return &resp

	} else if tenantType == string(SuperTenant) {

		nst, responseCode, err := createSuperTenant(tenantAlias, subdomain, superServicesConfig, publicServicesConfig, privateAccessConfig, customAccessConfig, liegeUUID, createdBy)
		if err != nil {
			resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
			return &resp
		}
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
				is_super_tenant,
				liege_uuid,
				lord_uuid,
				created_by
					)   SELECT UUID_TO_BIN(?), ?, top_level_domain, secondary_domain, subscripify_deployment_cloud_location, ?, ?, UUID_TO_BIN(?),
				 UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), ?, UUID_TO_BIN(?), UUID_TO_BIN(?), ? FROM tenant WHERE tenant_uuid = UUID_to_bin(?);`

		_, rc, err := tenantdbserv.HttpResponseHelperSQLInsert(tdb.ExecContext(ctx, insertStr,
			nst.getTenantUUID(),
			nst.getAlias(),
			nst.getSubdomainName(),
			nst.getKubeNamespacePrefix(),
			nst.getSuperServicesConfig(),
			nst.getPublicServicesConfig(),
			nst.getPrivateAccessConfig(),
			nst.getCustomAccessConfig(),
			nst.isSuperTenant(),
			nst.getLiegeUUID(),
			nst.getLiegeUUID(),
			nst.getTenantCreator(),
			nst.getLiegeUUID()))

		if err != nil {
			resp.generateHttpResponseCodeAndMessage(rc, err.Error())
			return &resp
		}

		resp.generateHttpResponseCodeAndMessage(200, "created new super tenant")
		resp.generateNewTenantResponse(nst.getTenantUUID())

		return &resp

	} else if tenantType == string(LordTenant) {
		resp.generateHttpResponseCodeAndMessage(404, "lord tenants need to be created using the NewLordTenant function a POST to /lord-tenants")
		return &resp
	} else {
		resp.generateHttpResponseCodeAndMessage(404, "this does not seem to be a valid tenant type - only super or main tenants can be created with this endpoint")
		return &resp
	}

}

func GetTenant(
	tenantUUID string,
	creator string) iHttpResponse {
	var resp httpResponseData

	l, responseCode, err := loadOneTenantFromDatabase(tenantUUID, creator)

	if err != nil {
		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
		return &resp
	}
	resp.generateHttpResponseCodeAndMessage(200, "successful object sent")
	resp.generateLoadedTenantResponse(l, GET)

	return &resp
}

func DeleteTenant(tenantUUID string,
	creator string) iHttpResponse {
	var resp httpResponseData

	//this will take care of validating inputs and loading the object
	l, responseCode, err := loadOneTenantFromDatabase(tenantUUID, creator)

	if err != nil {
		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
		return &resp
	}
	deleteStr := `DELETE FROM tenant WHERE tenant_uuid = UUID_TO_BIN(?)`
	//setting up a 10 second timeout (could be less)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tdb := tenantdbserv.Tdb.Handle
	_, rc, err := tenantdbserv.HttpResponseHelperSQLDelete(tdb.ExecContext(ctx, deleteStr, l.tenantUUID.String()))
	if err != nil {
		resp.generateHttpResponseCodeAndMessage(rc, err.Error())
		return &resp
	}
	resp.generateHttpResponseCodeAndMessage(200, "successful object removed")
	resp.generateLoadedTenantResponse(l, DELETE)

	return &resp
}

func UpdateTenant(
	tenantUUID string,
	tenantAlias string,
	lordServicesConfig string,
	superServicesConfig string,
	publicServicesConfig string,
	privateAccessConfig string,
	customAccessConfig string,
	creator string) iHttpResponse {
	var resp httpResponseData

	//this will take care of validating inputs and loading the object
	l, responseCode, err := loadOneTenantFromDatabase(tenantUUID, creator)
	if err != nil {
		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
		return &resp
	}

	var setString strings.Builder
	fmt.Fprintf(&setString, "UPDATE tenant SET")
	if tenantAlias != "" {
		if err := l.setAlias(tenantAlias); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " tenant_alias = '%s',", tenantAlias)
		}
	}
	if lordServicesConfig != "" {
		if err := l.setLordServicesConfig(lordServicesConfig); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " lord_services_config = UUID_TO_BIN('%s'),", l.getLordServicesConfig())
		}
	}
	if superServicesConfig != "" {
		if err := l.setSuperServicesConfig(superServicesConfig); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " super_services_config = UUID_TO_BIN('%s'),", l.getSuperServicesConfig())
		}
	}
	if publicServicesConfig != "" {
		if err := l.setPublicServicesConfig(publicServicesConfig); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " public_services_config = UUID_TO_BIN('%s'),", l.getPublicServicesConfig())
		}
	}
	if privateAccessConfig != "" {
		if err := l.setPrivateAccessConfig(privateAccessConfig); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " private_access_config = UUID_TO_BIN('%s'),", l.getPrivateAccessConfig())
		}
	}
	if customAccessConfig != "" {
		if err := l.setCustomAccessConfig(customAccessConfig); err != nil {
			resp.generateHttpResponseCodeAndMessage(400, err.Error())
			return &resp
		} else {
			fmt.Fprintf(&setString, " custom_access_config = UUID_TO_BIN('%s'),", l.getCustomAccessConfig())
		}
	}

	var updateString strings.Builder

	fmt.Fprint(&updateString, strings.TrimSuffix(setString.String(), ","), " WHERE tenant_uuid = UUID_TO_BIN(?)")

	//setting up a 10 second timeout (could be less)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tdb := tenantdbserv.Tdb.Handle
	_, rc, err := tenantdbserv.HttpResponseHelperSQLUpdate(tdb.ExecContext(ctx, updateString.String(), l.tenantUUID.String()))
	if err != nil {
		resp.generateHttpResponseCodeAndMessage(rc, err.Error())
		return &resp
	}
	resp.generateHttpResponseCodeAndMessage(200, "successful object updated")
	resp.generateLoadedTenantResponse(l, PATCH)

	return &resp
}

func ListTenants(
	page int,
	perPage int,
	pipedTenantType string,
	pipedTenantUUIDs string,
	pipedTenantAliases string,
	pipedSubdomains string,
	pipedDomains string,
	pipedConfigUUIDs string,
	pipedConfigAliases string,
	pipedAccessUUIDs string,
	pipedAccessAliases string) iHttpResponse {
	var resp httpResponseData

	so, responseCode, err := createTenantSearchObject(pipedTenantType, pipedTenantUUIDs, pipedTenantAliases,
		pipedSubdomains, pipedDomains, pipedConfigUUIDs, pipedConfigAliases, pipedAccessUUIDs, pipedAccessAliases)

	if err != nil {
		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
		return &resp
	}

	sr, responseCode, err := so.buildTenantSearchResults(page, perPage)
	if err != nil {
		resp.generateHttpResponseCodeAndMessage(responseCode, err.Error())
		return &resp
	}

	resp.generateHttpResponseCodeAndMessage(200, "successful")
	resp.generateTenantSearchList(sr, GET)

	return &resp
}
