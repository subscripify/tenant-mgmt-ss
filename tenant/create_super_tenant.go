package tenant

import (
	"fmt"
)

type superTenant struct {
	tenant
}

// Define a Stringer interface that gives a string representation of the type
func (t superTenant) String() string {
	return fmt.Sprintf("This is a organization named %s", t.orgName)
}

func createSuperTenant(orgName string, subdomain string, createdBy string) iTenant {
	return &superTenant{
		tenant: tenant{
			orgName:             orgName,
			subdomain:           subdomain,
			createdBy:           createdBy,
			kubeNamespacePrefix: "newKube",

			tenantType: SuperTenant,
		},
	}
}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
