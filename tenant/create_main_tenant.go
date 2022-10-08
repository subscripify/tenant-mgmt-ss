package tenant

import (
	"fmt"
)

type mainTenant struct {
	tenant
}

// Define a Stringer interface that gives a string representation of the type
func (t mainTenant) String() string {
	return fmt.Sprintf("This is a organization named %s", t.alias)
}

func createMainTenant(
	tenantAlias string,
	subdomain string,
	publicServicesConfig string,
	customAccessConfig string,
	liegeUUID string,
	createdBy string) (iTenant, int, error) {

	var m mainTenant

	//all of these checks would fail due to first level validation and should be 400s
	if err := m.setTenantType(MainTenant); err != nil {

		return nil, 400, err
	}
	if err := m.setAlias(tenantAlias); err != nil {

		return nil, 400, err
	}
	if err := m.setSubdomainName(subdomain); err != nil {

		return nil, 400, err
	}
	//there is no Lord Services config or Super services config for main tenants
	if err := m.setPublicServicesConfig(publicServicesConfig); err != nil {

		return nil, 400, err
	}
	//there is no Private Access config for main tenants only custom configs
	if err := m.setCustomAccessConfig(customAccessConfig); err != nil {

		return nil, 400, err
	}
	if err := m.setCreatedBy(createdBy); err != nil {

		return nil, 400, err
	}
	if err := m.setLiegeUUID(liegeUUID); err != nil {

		return nil, 400, err
	}

	// if all pass - set the new tenant id
	m.setNewTenantUUID()
	return &m, 200, nil
}
