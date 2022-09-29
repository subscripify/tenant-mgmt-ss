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
	createdBy string) (iTenant, int, error) {

	var s superTenant

	//all of these checks would fail due to first level validation and should be 400s
	if err := s.setTenantType(SuperTenant); err != nil {

		return nil, 400, err
	}
	if err := s.setAlias(tenantAlias); err != nil {

		return nil, 400, err
	}
	if err := s.setSubdomainName(subdomain); err != nil {

		return nil, 400, err
	}

	if err := s.setSuperServicesConfig(superServicesConfig); err != nil {

		return nil, 400, err
	}

	if err := s.setPublicServicesConfig(publicServicesConfig); err != nil {

		return nil, 400, err
	}
	//there is no Private Access config for main tenants only custom configs
	if err := s.setPrivateAccessConfig(privateAccessConfig); err != nil {

		return nil, 400, err
	}
	if err := s.setCustomAccessConfig(customAccessConfig); err != nil {

		return nil, 400, err
	}
	if err := s.setCreatedBy(createdBy); err != nil {

		return nil, 400, err
	}
	if err := s.setLiegeUUID(liegeUUID); err != nil {

		return nil, 400, err
	}

	// if all pass - set the new tenant id
	s.setNewTenantUUID()
	return &s, 200, nil

}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
