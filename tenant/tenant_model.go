package tenant

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type iTenant interface {
	setNewTenantUUID() error
	setOrgName(orgName string)
	setSubdomainName(subdomain string) error
	setInternalServicesConfig(internalServicesConfig uuid.UUID) error
	setSuperServicesConfig(superServicesConfig uuid.UUID) error
	setPublicServicesConfig(publicServicesConfig uuid.UUID) error
	setPrivateServicesConfig(publicServicesConfig uuid.UUID) error
	setCustomServicesConfig(publicServicesConfig uuid.UUID) error
	setCloudLocation(cloudLocation CloudLocation) error
	setTenantType(tenantType TenantType) error
	setCreatedBy(userIdentifier string) error
	setCreateTimestamp()
	GetTenant() (string, error)
}

type tenant struct {
	tenantUUID             uuid.UUID
	orgName                string
	subdomain              string
	secondaryDomain        string
	topLevelDomain         string
	kubeNamespacePrefix    string
	internalServicesConfig uuid.UUID
	superServicesConfig    uuid.UUID
	publicServicesConfig   uuid.UUID
	privateServicesConfig  uuid.UUID
	customServicesConfig   uuid.UUID
	cloudLocation          CloudLocation
	leigeUUID              uuid.UUID
	isLordTenant           bool
	isSuperTenant          bool
	createTimestamp        time.Time
	createdBy              string
}

type TenantType string

const (
	LordTenant  TenantType = "lord"
	SuperTenant            = "super"
	MainTenant             = "main"
)

type CloudLocation string

const (
	Azure CloudLocation = "Azure"
	ACS                 = "ACS"
	GCP                 = "GCP"
)

func (t *tenant) setNewTenantUUID() error {
	t.tenantUUID = uuid.New() //this is the only place in the application this value is created - but still check for duplicate - if error loop again - if loop three times return error
	return nil
}

// the org name does not need to be unique and is alias used for quick reference when searching
func (t *tenant) setOrgName(orgName string) {
	t.orgName = orgName
}

// the subdomain name must be unique within the same lord tenant. A lord tenant ties to a second level domain (eg. subscripify.com)
func (t *tenant) setSubdomainName(subdomain string) error {
	pattern := "^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$"
	r := regexp.MustCompile(pattern)
	if r.MatchString(subdomain) {
		t.subdomain = subdomain
		t.kubeNamespacePrefix = subdomain
		return nil
	}
	err := fmt.Errorf("Not a valid subdomain name - must match pattern '^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$'")
	return err
}

// all lord tenants require this config other types of tenants this is null - it refers to the configuration files for internal services available to this tenant
func (t *tenant) setInternalServicesConfig(internalServicesConfig uuid.UUID) error {

	if !t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting internal services - lord tenants only")
		return err
	}
	t.internalServicesConfig = internalServicesConfig
	return nil
}

// all lord tenants and super tenants require this config other types of tenants this is null - it refers to the configuration files for super services available to this tenant
func (t *tenant) setSuperServicesConfig(superServicesConfig uuid.UUID) error {

	if !t.isLordTenant && !t.isSuperTenant {
		err := fmt.Errorf("invalid tenant type for setting super services - lord tenants and super tenants only")
		return err
	}
	t.superServicesConfig = superServicesConfig
	return nil
}

// all tenants require this config - it refers to the configuration files for public services available to this tenant
func (t *tenant) setPublicServicesConfig(publicServicesConfig uuid.UUID) error {

	t.publicServicesConfig = publicServicesConfig
	return nil
}

// all all Super Tenants require this config and is set at time of creation - it refers to the private services oAuth/access configuration for the super tenant
func (t *tenant) setPrivateServicesConfig(publicServicesConfig uuid.UUID) error {

	if !t.isSuperTenant {
		err := fmt.Errorf("invalid tenant type for setting Private Services Config - super tenants only")
		return err
	}
	t.publicServicesConfig = publicServicesConfig
	return nil
}

// all all super and main tenants require this config and is set at time of creation - it refers to the custom services api oAuth/access configuration for this tenant
func (t *tenant) setCustomServicesConfig(publicServicesConfig uuid.UUID) error {

	if t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting Private Services Config - super tenants and main tenants only")
		return err
	}
	t.publicServicesConfig = publicServicesConfig
	return nil
}

// limited to Azure, ACS, GCP
func (t *tenant) setCloudLocation(cloudLocation CloudLocation) error {
	t.cloudLocation = cloudLocation
	return nil
}

// sets the tenant type booleans in the database based upon the passed value
func (t *tenant) setTenantType(tenantType TenantType) error {

	switch tenantType {
	case LordTenant:
		t.isLordTenant = true
		t.isSuperTenant = false
	case SuperTenant:
		t.isLordTenant = false
		t.isSuperTenant = true
	case MainTenant:
		t.isLordTenant = false
		t.isSuperTenant = false
	default:
		return fmt.Errorf("invalid tenant type")
	}
	return nil
}

// id from identity system - must validate as email of current user session
func (t *tenant) setCreatedBy(userIdentifier string) error {
	t.createdBy = userIdentifier
	return nil
}

func (t *tenant) setCreateTimestamp() {
	t.createTimestamp = time.Now()
}

// return a raw json of the current tenant otherwise return error
func (t *tenant) GetTenant() (string, error) {

	return "a json ", nil
}
