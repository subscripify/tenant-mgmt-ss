package tenant

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/tenantdb"
	"github.com/google/uuid"
)

func NewLordTenant(
	tenantAlias string,
	topLevelDomain string,
	secondaryDomain string,
	subdomain string,
	lordServicesConfig string,
	superServicesConfig string,
	publicServicesConfig string,
	cloudLocation CloudLocation,
	createdBy string) (uuid.UUID, error) {
	//no special processing required - this is a pass through to maintain the pattern. the NewTenant function covers the factory for non lord tenant types
	nlt, err := createLordTenant(tenantAlias, topLevelDomain, secondaryDomain, subdomain, lordServicesConfig, superServicesConfig, publicServicesConfig, cloudLocation, createdBy)
	if err != nil {
		log.Printf("new lord tenant object creation fail: %v", err)
		return uuid.UUID{}, err
	}

	//setting up a 10 second timeout (could be less)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tdb := tenantdbserv.Tdb.Handle

	if err != nil {
		log.Printf("new lord tenant object get fail: %v", err)

		return uuid.UUID{}, err
	}
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

	r, err := tdb.ExecContext(ctx, insertStr,
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
		nlt.getTenantCreator())

	if err != nil {
		log.Printf("fail on insert: %v", err)
		return uuid.UUID{}, err
	}

	count, _ := r.RowsAffected()
	log.Printf("number of rows inserted :%v", count)
	return nlt.getTenantUUID(), nil
}

func NewTenant(
	tenantType TenantType,
	tenantAlias string,
	subdomain string,
	superServicesConfig string,
	publicServicesConfig string,
	privateAccessConfig string,
	customAccessConfig string,
	liegeUUID string,
	createdBy string) (iTenant, error) {
	// factory function for new super or main tenants. this function will also look up lord tenants using liege tenant value

	//setting up a 10 second timeout (could be less)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var lordUUID string
	var isSupperTenant bool = false
	var isLordTenant bool = false
	t := tenantdbserv.Tdb.Handle
	if tenantType == MainTenant {
		err := t.QueryRowContext(ctx, "SELECT lord_uuid, is_super_tenant from tenants where tenant_uuid = ?", liegeUUID).Scan(&lordUUID, &isSupperTenant)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("could not find the liege tenant specified")
		}
		if !isSupperTenant {
			return nil, fmt.Errorf("this not a valid liege tenant for a main tenant, must use a super tenant")
		}
		return createMainTenant(tenantAlias, subdomain, publicServicesConfig, customAccessConfig, liegeUUID, lordUUID, createdBy), nil
	} else if tenantType == SuperTenant {
		err := t.QueryRowContext(ctx, "SELECT is_lord_tenant from tenants where tenant_uuid = ?", liegeUUID).Scan(&isLordTenant)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("could not find the liege tenant specified")
		}
		if !isLordTenant {
			return nil, fmt.Errorf("this not a valid liege tenant for a super tenant, must use a lord tenant")
		}
		return createSuperTenant(tenantAlias, subdomain, superServicesConfig, publicServicesConfig, privateAccessConfig, customAccessConfig, liegeUUID, lordUUID, createdBy), nil
	} else if tenantType == LordTenant {
		err := fmt.Errorf("lord tenants need to be created using the NewLordTenant function")
		return nil, err
	}

	return nil, fmt.Errorf("no such tenant type")
}
