package tenant

import (
	"fmt"
)

type lordTenant struct {
	tenant
}

// Define a Stringer interface that gives a string representation of the type
func (t lordTenant) String() string {
	return fmt.Sprintf("This is a organization named %s", t.orgName)
}

func createLordTenant(orgName string, subdomain string, createdBy string) iTenant {
	return &mainTenant{
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
