package tenant

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/subscripify-logger"
	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/tenantdb"
	"github.com/google/uuid"
)

type loadedTenant struct {
	tenant
}

type loader struct {
	tenantUUID           sql.NullString
	alias                sql.NullString
	subdomain            sql.NullString
	secondaryDomain      sql.NullString
	topLevelDomain       sql.NullString
	kubeNamespacePrefix  sql.NullString
	lordServicesConfig   sql.NullString
	superServicesConfig  sql.NullString
	publicServicesConfig sql.NullString
	privateAccessConfig  sql.NullString
	customAccessConfig   sql.NullString
	cloudLocation        sql.NullString
	liegeUUID            sql.NullString
	lordUUID             sql.NullString
	lordTenant           sql.NullBool
	superTenant          sql.NullBool
	createTimestamp      sql.NullTime
	createdBy            sql.NullString
}

func loadOneTenantFromDatabase(tenantUUID string, creator string) (*loadedTenant, int, error) {
	// check ownership - return error if not owned
	// logic here would also check session  when implementing oAuth

	// get the tenant record based upon uuid
	var l loader
	var t loadedTenant

	tenantUUIDParsed, err := uuid.Parse(tenantUUID)
	if err != nil {
		return nil, 400, fmt.Errorf("failed to parse UUID in URL - is it a valid uuid?")
	}

	tdb := tenantdbserv.Tdb.Handle
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	selectStr := `SELECT 
	BIN_TO_UUID(tenant_uuid),
	 tenant_alias,
	 top_level_domain,
	 secondary_domain,
	 subdomain,
	 kube_namespace_prefix,
	 BIN_TO_UUID(lord_services_config),
	 BIN_TO_UUID(super_services_config),
	 BIN_TO_UUID(public_services_config),
	 BIN_TO_UUID(private_access_config),
	 BIN_TO_UUID(custom_access_config),
	 subscripify_deployment_cloud_location,
	 BIN_TO_UUID(liege_uuid),
	 BIN_TO_UUID(lord_uuid),
	 is_lord_tenant,
	 is_super_tenant,
	 create_timestamp,
	 created_by
	 FROM tenant WHERE tenant_uuid = UUID_TO_BIN(?);`

	r := tdb.QueryRowContext(ctx, selectStr, tenantUUIDParsed.String())
	if err := r.Scan(
		&l.tenantUUID,
		&l.alias,
		&l.topLevelDomain,
		&l.secondaryDomain,
		&l.subdomain,
		&l.kubeNamespacePrefix,
		&l.lordServicesConfig,
		&l.superServicesConfig,
		&l.publicServicesConfig,
		&l.privateAccessConfig,
		&l.customAccessConfig,
		&l.cloudLocation,
		&l.liegeUUID,
		&l.lordUUID,
		&l.lordTenant,
		&l.superTenant,
		&l.createTimestamp,
		&l.createdBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, fmt.Errorf("tenant uuid does not exist")
		} else {
			var re = regexp.MustCompile(`.+Error.1411.+for.function.uuid_to_bin`)

			if len(re.FindStringIndex(err.Error())) == 0 {
				return nil, 400, fmt.Errorf("failed to parse UUID in url at query level - is it a valid uuid?")
			}
			return nil, 500, fmt.Errorf("error reading database: %s", err.Error())
		}

	}

	if l.tenantUUID.Valid {
		id, err := uuid.Parse(l.tenantUUID.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		t.tenantUUID = id

	}
	if l.alias.Valid {
		t.alias = l.alias.String
	}

	if l.topLevelDomain.Valid {
		t.topLevelDomain = l.topLevelDomain.String
	}

	if l.secondaryDomain.Valid {
		t.secondaryDomain = l.secondaryDomain.String
	}

	if l.subdomain.Valid {
		t.subdomain = l.subdomain.String
	}

	if l.kubeNamespacePrefix.Valid {
		t.kubeNamespacePrefix = l.kubeNamespacePrefix.String
	}
	// start services configurations loading

	if l.lordServicesConfig.Valid {
		configID, err := uuid.Parse(l.lordServicesConfig.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		if !(l.lordTenant.Valid && l.lordTenant.Bool) {
			subscripifylogger.WarningLog.Printf("a lord_services_config uuid was loaded from db for a non-lord tenant")
		}
		t.lordServicesConfig = configID
	}

	if l.superServicesConfig.Valid {
		configID, err := uuid.Parse(l.superServicesConfig.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		if !((l.superTenant.Valid && l.superTenant.Bool) || (l.lordTenant.Valid && l.lordTenant.Bool)) {
			subscripifylogger.WarningLog.Printf("a super_tenant_config uuid was loaded from db for a non-super/non-lord tenant")
		}
		t.superServicesConfig = configID
	}

	if l.publicServicesConfig.Valid {
		configID, err := uuid.Parse(l.publicServicesConfig.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		t.publicServicesConfig = configID
	}

	//end services configuration loading
	// start access configuration loading

	if l.privateAccessConfig.Valid {
		configID, err := uuid.Parse(l.privateAccessConfig.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}

		if !(l.superTenant.Valid && l.superTenant.Bool) {
			subscripifylogger.WarningLog.Printf("a private_access_config uuid was loaded from db for a non-super tenant")
		}

		t.privateAccessConfig = configID
	}

	if l.customAccessConfig.Valid {
		configID, err := uuid.Parse(l.customAccessConfig.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}

		if !(l.superTenant.Valid && l.superTenant.Bool) {
			subscripifylogger.WarningLog.Printf("a custom_access_config uuid was loaded from db for a lord tenant")
		}

		t.customAccessConfig = configID
	}

	//end access configuration loading

	if l.cloudLocation.Valid {
		switch {
		case l.cloudLocation.String == string(Azure):
			t.cloudLocation = Azure
		case l.cloudLocation.String == string(GCP):
			t.cloudLocation = GCP
		case l.cloudLocation.String == string(AWS):
			t.cloudLocation = AWS
		default:
			return nil, 500, fmt.Errorf("error loading object: %s", "unsupported cloud provider in DB.subscripify_cloud_location")

		}

	}

	if l.liegeUUID.Valid {
		id, err := uuid.Parse(l.liegeUUID.String)
		if err != nil {
			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		t.liegeUUID = id
	}

	if l.lordUUID.Valid {
		id, err := uuid.Parse(l.lordUUID.String)
		if err != nil {

			return nil, 500, fmt.Errorf("error loading object: %s", err.Error())
		}
		t.lordUUID = id
	}

	// setting the tenant types

	if l.lordTenant.Valid && l.lordTenant.Bool {
		t.lordTenant = l.lordTenant.Bool

	} else {
		t.lordTenant = false
	}

	if l.superTenant.Valid && l.superTenant.Bool {
		t.superTenant = l.superTenant.Bool

	} else {
		t.superTenant = false
	}

	//end setting the tenant types

	if l.createTimestamp.Valid {
		t.createTimestamp = l.createTimestamp.Time
	}

	if l.createdBy.Valid {
		t.createdBy = l.createdBy.String

	}

	return &t, 200, nil

}
