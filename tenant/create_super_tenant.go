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
	createdBy string) (iTenant, iHttpResponse) {

	var s superTenant
	var r httpResponseData

	//all of these checks would fail due to first level validation and should be 400s
	if err := s.setTenantType(MainTenant); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setTenantType")
		return nil, &r
	}
	if err := s.setAlias(tenantAlias); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setAlias")
		return nil, &r
	}
	if err := s.setSubdomainName(subdomain); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSubdomainName")
		return nil, &r
	}

	if err := s.setSuperServicesConfig(superServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSuperServicesConfig")
		return nil, &r
	}

	if err := s.setPublicServicesConfig(publicServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setPublicServicesConfig")
		return nil, &r
	}
	//there is no Private Access config for main tenants only custom configs
	if err := s.setPrivateAccessConfig(privateAccessConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setPrivateAccessConfig")
		return nil, &r
	}
	if err := s.setCustomAccessConfig(customAccessConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCustomAccessConfig")
		return nil, &r
	}
	if err := s.setCreatedBy(createdBy); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCreatedBy")
		return nil, &r
	}
	if err := s.setLiegeUUID(liegeUUID); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setLiegeUUID")
		return nil, &r
	}

	// if all pass - set the new tenant id
	s.setNewTenantUUID()
	return &s, &r

}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
