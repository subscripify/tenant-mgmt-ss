package tenant

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
)

type tenantTestData struct {
	alias          string
	tenantTld      string
	tenantSecD     string
	lordTenantSub  string
	superTenantSub string
	mainTenantSub  string
	tenantCloud    CloudLocation
	createdBy      string
	privateAccess  uuid.UUID
	customAccess   uuid.UUID
	lordConfig     uuid.UUID
	superConfig    uuid.UUID
	publicConfig   uuid.UUID
}

// creates a new test data struct with everything that one needs to do a tenant test
func NewTenantTestData() *tenantTestData {
	t := tenantTestData{}
	t.alias = "testing-cIXzsBP4bw"
	t.tenantTld = "bz"
	t.tenantSecD = "test-domain"
	t.lordTenantSub = "lord-tenant"
	t.superTenantSub = "super-tenant"
	t.mainTenantSub = "main-tenant"
	t.tenantCloud = Azure
	t.createdBy = "william.ohara@subscripify.com"

	// these data pints are in the DB, need to develop the rest of the endpoints for config management
	t.privateAccess = uuid.MustParse("c9057cff-3863-11ed-907f-f5001f9bae96")
	t.customAccess = uuid.MustParse("845c5fe4-3864-11ed-907f-f5001f9bae96")
	t.lordConfig = uuid.MustParse("60135d7c-3857-11ed-907f-f5001f9bae96")
	t.superConfig = uuid.MustParse("432ead01-385a-11ed-907f-f5001f9bae96")
	t.publicConfig = uuid.MustParse("6a24689f-385a-11ed-907f-f5001f9bae96")

	return &t

}

func TestDataCreate() {
	testD := NewTenantTestData()

	// outer loop creates three lord tenants
	for i := 0; i < 4; i++ {
		seedl := time.Now().UTC().UnixNano()
		nameGeneratorl := namegenerator.NewNameGenerator(seedl)
		ltname := nameGeneratorl.Generate()
		nlt := NewLordTenant((ltname + "-ta"), "com", (ltname + "-d"), "lordtenant", testD.lordConfig.String(),
			testD.superConfig.String(), testD.publicConfig.String(), testD.tenantCloud, testD.createdBy)
		nltId := nlt.GetHttpResponse().NewTenant.TenantUUID
		log.Printf("created lord tenant: %s", nltId)
		// mid loop creates 10 super tenants per lord
		for j := 0; j < 11; j++ {
			seeds := time.Now().UTC().UnixNano()
			nameGenerators := namegenerator.NewNameGenerator(seeds)
			stname := nameGenerators.Generate()
			nst := NewTenant("super", (stname + "-ta"), ("st-" + stname + fmt.Sprint(i+j)), testD.superConfig.String(),
				testD.publicConfig.String(), testD.privateAccess.String(), testD.customAccess.String(), nltId, "william.ohara@subscripify.com")
			nstId := nst.GetHttpResponse().NewTenant.TenantUUID
			log.Printf("created super tenant: %s", nstId)
			// inner loop creates 30 main tenants per super
			for k := 0; k < 31; k++ {
				seedm := time.Now().UTC().UnixNano()
				nameGeneratorm := namegenerator.NewNameGenerator(seedm)
				mtname := nameGeneratorm.Generate()
				nst := NewTenant("main", (stname + "-ta"), ("mt-" + mtname + fmt.Sprint(i+j+k)), testD.superConfig.String(),
					testD.publicConfig.String(), testD.privateAccess.String(), testD.customAccess.String(), nstId, "william.ohara@subscripify.com")
				nmtId := nst.GetHttpResponse().NewTenant.TenantUUID
				log.Printf("created main tenant: %s", nmtId)
			}
		}
	}
}
