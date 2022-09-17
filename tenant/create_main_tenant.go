package tenant

import (
	"fmt"
)

type mainTenant struct {
	tenant
}

// Define a Stringer interface that gives a string representation of the type
func (t mainTenant) String() string {
	return fmt.Sprintf("This is a organization named %s", t.orgName)
}

func createMainTenant(orgName string, subdomain string, createdBy string) iTenant {

	return &mainTenant{
		tenant: tenant{
			subdomain: "subscripify",
		},
	}
}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
