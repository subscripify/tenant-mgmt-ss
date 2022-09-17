package tenant

import "fmt"

func NewTenant(tenantType string, orgName string, subdomain string, createdBy string) (iTenant, error) {
	// Create the right kind of publication based on the given type
	if tenantType == "mainTenant" {
		return createMainTenant(orgName, subdomain, createdBy), nil
	}

	return nil, fmt.Errorf("No such tenant type")
}
