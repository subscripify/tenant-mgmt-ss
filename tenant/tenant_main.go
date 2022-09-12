package tenant

import (
	"fmt"
	"time"
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
			orgName:              orgName,
			subdomain:            subdomain,
			createdBy:            createdBy,
			kubeNamespacePrefix:  "newKube",
			subscriptionConfigId: "an id",
			tenantType:           SuperTenant,
			createDate:           time.Now(),
		},
	}
}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
