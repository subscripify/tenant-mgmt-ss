package tenant

type lordTenant struct {
	tenant
}

func createLordTenant(
	tenantAlias string,
	topLevelDomain string,
	secondaryDomain string,
	subdomain string,
	lordServicesConfig string,
	superServicesConfig string,
	publicServicesConfig string,
	cloudLocation CloudLocation,
	createdBy string) (iTenant, iHttpResponse) {
	var l lordTenant
	var r httpResponseData
	//all of these are due to first level validation failures and should be 400s
	if err := l.setTenantType(LordTenant); err != nil {

		r.logAndGenerateHttpResponseData(400, err.Error(), "setTenantType")
		return nil, &r
	}
	if err := l.setSubdomainName(subdomain); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSubdomainName")
		return nil, &r
	}
	if err := l.setSecondaryDomainName(secondaryDomain); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSecondaryDomainName")
		return nil, &r
	}
	if err := l.setTopLevelDomain(topLevelDomain); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setTopLevelDomain")
		return nil, &r
	}
	if err := l.setSuperServicesConfig(superServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setSuperServicesConfig")
		return nil, &r
	}
	if err := l.setPublicServicesConfig(publicServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setPublicServicesConfig")
		return nil, &r
	}
	if err := l.setLordServicesConfig(lordServicesConfig); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setLordServicesConfig")
		return nil, &r
	}
	if err := l.setCloudLocation(cloudLocation); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCloudLocation")
		return nil, &r
	}
	if err := l.setCreatedBy(createdBy); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setCreatedBy")
		return nil, &r
	}
	if err := l.setAlias(tenantAlias); err != nil {
		r.logAndGenerateHttpResponseData(400, err.Error(), "setAlias")
		return nil, &r
	}

	l.setNewTenantUUID()

	return &l, &r

}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
