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
	createdBy string) (iTenant, error) {
	var l lordTenant

	//parse the UUID strings to ensure that they are UUIDs

	if err := l.setTenantType(LordTenant); err != nil {
		return nil, err
	}
	if err := l.setSubdomainName(subdomain); err != nil {
		return nil, err
	}
	if err := l.setSecondaryDomainName(secondaryDomain); err != nil {
		return nil, err
	}
	if err := l.setTopLevelDomain(topLevelDomain); err != nil {
		return nil, err
	}
	if err := l.setLordServicesConfig(lordServicesConfig); err != nil {
		return nil, err
	}
	if err := l.setSuperServicesConfig(superServicesConfig); err != nil {
		return nil, err
	}
	if err := l.setPublicServicesConfig(publicServicesConfig); err != nil {
		return nil, err
	}
	if err := l.setCloudLocation(cloudLocation); err != nil {
		return nil, err
	}
	if err := l.setCreatedBy(createdBy); err != nil {
		return nil, err
	}
	if err := l.setAlias(tenantAlias); err != nil {
		return nil, err
	}

	l.setNewTenantUUID()

	return &l, nil

}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
