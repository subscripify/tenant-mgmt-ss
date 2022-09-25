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
	createdBy string) (iTenant, iHttpResponse) {

	var m mainTenant
	var r httpResponseData

	//all of these checks would fail due to first level validation and should be 400s
	if err := m.setTenantType(MainTenant); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setTenantType")
		return nil, &r
	}
	if err := m.setAlias(tenantAlias); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setAlias")
		return nil, &r
	}
	if err := m.setSubdomainName(subdomain); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSubdomainName")
		return nil, &r
	}
	//there is no Lord Services config or Super services config for main tenants
	if err := m.setPublicServicesConfig(publicServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setPublicServicesConfig")
		return nil, &r
	}
	//there is no Private Access config for main tenants only custom configs
	if err := m.setCustomAccessConfig(customAccessConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCustomAccessConfig")
		return nil, &r
	}
	if err := m.setCreatedBy(createdBy); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCreatedBy")
		return nil, &r
	}
	if err := m.setLiegeUUID(liegeUUID); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setLiegeUUID")
		return nil, &r
	}

	// if all pass - set the new tenant id
	m.setNewTenantUUID()
	return &m, &r
}
