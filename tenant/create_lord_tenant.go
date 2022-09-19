package tenant

type lordTenant struct {
	tenant
}

func createLordTenant(
	orgName string,
	subdomain string,
	secondaryDomain string,
	topLevelDomain string,
	internalServicesConfigAlias string,
	superServicesConfigAlias string,
	publicServicesConfigAlias string,
	cloudLocation CloudLocation,
	createdBy string) (iTenant, error) {
	var l lordTenant

	newUUID := l.setNewTenantUUID()
	l.setOrgName(orgName)
	if err := l.setSubdomainName(subdomain); err != nil {
		return nil, err
	}
	if err := l.setSecondaryDomainName(secondaryDomain); err != nil {
		return nil, err
	}
	if err := l.setTopLevelDomain(topLevelDomain); err != nil {
		return nil, err
	}
	if err := l.setTenantType(LordTenant); err != nil {
		return nil, err
	}
	if err := l.setInternalServicesConfig(internalServicesConfigAlias); err != nil {
		return nil, err
	}

	return &l, nil
	// return &lordTenant{
	// 	tenant: tenant{
	// 		orgName:             orgName,
	// 		subdomain:           subdomain,
	// 		createdBy:           createdBy,
	// 		kubeNamespacePrefix: "newKube",
	// 		tenantType: SuperTenant,
	// 	},
	// }
}

//createSuperTenant
//find Tenant
//find tenants by owner org
//delete tenant
//change tenant
