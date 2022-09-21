package tenant

import (
	"fmt"
)

type superTenant struct {
	tenant
}

// Define a Stringer interface that gives a string representation of the type
func (t superTenant) String() string {
	return fmt.Sprintf("This is a organization named %s", t.alias)
}

func createSuperTenant(
	tenantAlias string,
	subdomain string,
	superServicesConfig string,
	publicServicesConfig string,
	privateAccessConfig string,
	customAccessConfig string,
	liegeUUID string,
	lordUUID string,
	createdBy string) iTenant {
	return &superTenant{
		tenant: tenant{
			alias:               tenantAlias,
			subdomain:           subdomain,
			createdBy:           createdBy,
			kubeNamespacePrefix: "newKube",
		},
	}
}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
