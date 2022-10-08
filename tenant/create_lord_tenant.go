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
	createdBy string) (iTenant, int, error) {
	var l lordTenant

	//all of these are due to first level validation failures and should be 400s
	if err := l.setTenantType(LordTenant); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setSubdomainName(subdomain); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setSecondaryDomainName(secondaryDomain); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setTopLevelDomain(topLevelDomain); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setSuperServicesConfig(superServicesConfig); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setPublicServicesConfig(publicServicesConfig); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setLordServicesConfig(lordServicesConfig); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setCloudLocation(cloudLocation); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setCreatedBy(createdBy); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}
	if err := l.setAlias(tenantAlias); err != nil {

		return &lordTenant{tenant{}}, 400, err
	}

	l.setNewTenantUUID()

	return &l, 200, nil

}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
