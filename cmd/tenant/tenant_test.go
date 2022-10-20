package tenant

import (
	"log"
	"testing"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/tenant-mgmt-ss.git/internal/tenantdb"
)

// this function fills the database for testing search

func TestCreateTenantTree(t *testing.T) {
	// need to utilize stored procedures to protect database
	t.Cleanup(func() {
		testD := NewTenantTestData()
		tdb := tenantdbserv.Tdb.Handle
		_, err := tdb.Exec("DELETE FROM tenant WHERE tenant_alias = ? AND is_lord_tenant IS NULL AND is_super_tenant = false;", testD.alias)
		if err != nil {
			log.Printf("test cleanup failed on main tenant delete")
		}
		_, err = tdb.Exec("DELETE FROM tenant WHERE tenant_alias = ? AND is_lord_tenant IS NULL AND is_super_tenant = true;", testD.alias)
		if err != nil {
			log.Printf("test cleanup failed on super tenant delete")
		}
		_, err = tdb.Exec("DELETE FROM tenant WHERE tenant_alias = ? AND is_lord_tenant AND is_super_tenant = false;", testD.alias)
		if err != nil {
			log.Printf("test cleanup failed on lord tenant delete")
		}
	})

	testD := NewTenantTestData()
	tdb := tenantdbserv.Tdb.Handle

	nlt := NewLordTenant(testD.alias, testD.tenantTld, testD.tenantSecD,
		testD.lordTenantSub, testD.lordConfig.String(), testD.superConfig.String(),
		testD.publicConfig.String(), testD.tenantCloud, testD.createdBy)
	ltr := nlt.GetHttpResponse()

	goodLord := true
	clr := tdb.QueryRow("SELECT count(tenant_uuid) as count FROM tenant WHERE tenant_uuid = UUID_TO_BIN(?)", ltr.NewTenant.TenantUUID)
	lcount := 0
	_ = clr.Scan(&lcount)

	if ltr.HttpResponseCode != 200 {
		goodLord = false
	}
	if lcount != 1 {
		goodLord = false
	}

	if !goodLord {
		t.Logf("creation of lord tenant failed, %s", ltr.Message)
		t.Fail()
	} else {
		// use the lord tenant from the previous step
		nst := NewTenant("super", testD.alias, testD.superTenantSub, testD.superConfig.String(), testD.publicConfig.String(),
			testD.privateAccess.String(), testD.customAccess.String(), ltr.NewTenant.TenantUUID, testD.createdBy)
		str := nst.GetHttpResponse()

		goodSuper := true
		csr := tdb.QueryRow("SELECT count(tenant_uuid) as count FROM tenant WHERE tenant_uuid = UUID_TO_BIN(?)", str.NewTenant.TenantUUID)
		scount := 0
		_ = csr.Scan(&scount)

		if ltr.HttpResponseCode != 200 {
			goodSuper = false
		}
		if scount != 1 {
			goodSuper = false
		}

		if !goodSuper {
			t.Logf("creation of super tenant failed, %s", ltr.Message)
			t.Fail()
		} else {

			// use the super tenant from the previous step
			nmt := NewTenant("main", testD.alias, testD.mainTenantSub, "", testD.publicConfig.String(),
				"", testD.customAccess.String(), str.NewTenant.TenantUUID, testD.createdBy)
			mtr := nmt.GetHttpResponse()
			goodMain := true
			cmr := tdb.QueryRow("SELECT count(tenant_uuid) as count FROM tenant WHERE tenant_uuid = UUID_TO_BIN(?)", mtr.NewTenant.TenantUUID)
			mcount := 0
			_ = cmr.Scan(&mcount)

			if ltr.HttpResponseCode != 200 {
				goodMain = false
			}
			if mcount != 1 {
				goodMain = false
			}

			if !goodMain {
				t.Logf("creation of main tenant failed, %s", ltr.Message)
				t.Fail()
			}
		}
	}

}
