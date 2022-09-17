package tenant

type lordTenant struct {
	tenant
}

func createLordTenant(orgName string, subdomain string, createdBy string) iTenant {
	var l lordTenant
	l.setNewTenantUUID()
	l.setCreatedBy(createdBy)
	l.setSubdomainName(subdomain)

	return &l
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
