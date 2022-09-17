package tenant

type lordTenant struct {
	tenant
}

func createLordTenant(
	orgName string,
	subdomain string,
	internalServicesConfigAlias string,
	superServicesConfigAlias string,
	publicServicesConfigAlias string,
	cloudLocation CloudLocation,
	createdBy string) (iTenant, error) {
	var l lordTenant

	l.setNewTenantUUID()
	l.setTenantType(LordTenant)
	l.setCreatedBy(createdBy)
	l.setSubdomainName(subdomain)

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
