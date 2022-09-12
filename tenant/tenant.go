package tenant

import (
	"time"

	"github.com/google/uuid"
)

type iTenant interface {
	setNewOrgId() error
	setOrgName(orgName string)
	setSubdomainName(subdomain string) error
	setSubscriptionConfigId(subscriptionConfigId string)
	setTenantType(subscriptionConfigId string)
	setCreatedBy(userIdentifier string)
	setCreateDate()
}

type tenant struct {
	orgId                uuid.UUID
	orgName              string
	subdomain            string
	kubeNamespacePrefix  string
	subscriptionConfigId string
	tenantType           TenantType
	createdBy            string
	createDate           time.Time
}

type TenantType int

const (
	LordTenant TenantType = iota
	SuperTenant
)

func (t *tenant) setNewOrgId() error {
	t.orgId = uuid.New() //this is the only place in the application this value is created - but still check for duplicate - if error loop again - if loop three times return error
	return nil
}

func (t *tenant) setOrgName(orgName string) {
	t.orgName = orgName
}

func (t *tenant) setSubdomainName(subdomain string) error {
	t.subdomain = subdomain // check against database to make sure the name does not exist - no duplicates
	return nil
}

func (t *tenant) setSubscriptionConfigId(subscriptionConfigId string) {
	t.subscriptionConfigId = subscriptionConfigId
}

func (t *tenant) setTenantType(subscriptionConfigId string) {
	t.subscriptionConfigId = subscriptionConfigId
}

func (t *tenant) setCreatedBy(userIdentifier string) {
	t.createdBy = userIdentifier //id from identity system - must validate as email of current user session
}

func (t *tenant) setCreateDate() {
	t.createDate = time.Now()
}
