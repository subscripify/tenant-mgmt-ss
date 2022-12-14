package tenant

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bombsimon/tld-validator"
	"github.com/google/uuid"

	goaway "github.com/TwiN/go-away"
)

type iTenant interface {
	setNewTenantUUID()
	setTenantUUID(tenantUUID string) error
	setAlias(alias string) error
	setSubdomainName(subdomain string) error
	setSecondaryDomainName(secondaryDomain string) error
	setTopLevelDomain(topLevelDomain string) error
	setLordServicesConfig(lordServicesConfig string) error
	setSuperServicesConfig(superServicesConfig string) error
	setPublicServicesConfig(publicServicesConfig string) error
	setPrivateAccessConfig(privateAccessConfig string) error
	setCustomAccessConfig(customAccessConfig string) error
	setCloudLocation(cloudLocation CloudLocation) error
	setLiegeUUID(liegeUUID string) error
	setLordUUID(lordUUID string) error
	setTenantType(tenantType TenantType) error
	setCreatedBy(userIdentifier string) error
	getTenantUUID() uuid.UUID
	getAlias() string
	getSubdomainName() string
	getKubeNamespacePrefix() string
	getSecondaryDomainName() string
	getTopLevelDomain() string
	getLordServicesConfig() uuid.UUID
	getSuperServicesConfig() uuid.UUID
	getPublicServicesConfig() uuid.UUID
	getPrivateAccessConfig() uuid.UUID
	getCustomAccessConfig() uuid.UUID
	getCloudLocation() CloudLocation
	getLiegeUUID() uuid.UUID
	getLordUUID() uuid.UUID
	isLordTenant() bool
	isSuperTenant() bool
	getTenantCreator() string
	getCreateTime() time.Time
}

type tenant struct {
	tenantUUID           uuid.UUID     `db:"tenant.tenant_uuid"`
	alias                string        `db:"tenant.tenant_alias"`
	subdomain            string        `db:"tenant.subdomain"`
	secondaryDomain      string        `db:"tenant.secondary_domain"`
	topLevelDomain       string        `db:"tenant.top_level_domain"`
	kubeNamespacePrefix  string        `db:"tenant.kube_namespace_prefix"`
	lordServicesConfig   uuid.UUID     `db:"tenant.lord_services_config,omitempty"`
	superServicesConfig  uuid.UUID     `db:"tenant.super_services_config,omitempty"`
	publicServicesConfig uuid.UUID     `db:"tenant.public_services_config"`
	privateAccessConfig  uuid.UUID     `db:"tenant.private_access_config,omitempty"`
	customAccessConfig   uuid.UUID     `db:"tenant.custom_access_config,omitempty"`
	cloudLocation        CloudLocation `db:"tenant.subscripify_deployment_cloud_location"`
	liegeUUID            uuid.UUID     `db:"tenant.liege_uuid,omitempty"`
	lordUUID             uuid.UUID     `db:"tenant.lord_uuid,omitempty"`
	lordTenant           bool          `db:"tenant.is_lord_tenant,omitempty"`
	superTenant          bool          `db:"tenant.is_super_tenant,omitempty"`
	createTimestamp      time.Time     `db:"tenant.create_timestamp"`
	createdBy            string        `db:"tenant.created_by"`
}

type TenantType string

const (
	LordTenant  TenantType = "lord"
	SuperTenant TenantType = "super"
	MainTenant  TenantType = "main"
)

type CloudLocation string

const (
	Azure CloudLocation = "Azure"
	AWS   CloudLocation = "AWS"
	GCP   CloudLocation = "GCP"
)

// this is the only place in the application this value is created
func (t *tenant) setNewTenantUUID() {
	t.tenantUUID = uuid.New()
}

// this function loads a previously created UUID from a tenant record retrieved from the database
// this is used for loading a tenant object from db - not for creating a new tenant
func (t *tenant) setTenantUUID(tenantUUID string) error {

	loadedTenantUUID, err := uuid.Parse(tenantUUID)
	if err != nil {
		return fmt.Errorf("the tenant UUID failed to parse: %s", err)
	}
	t.tenantUUID = loadedTenantUUID
	return nil
}

// returns the byte[16] of the tenant uuid
func (t *tenant) getTenantUUID() uuid.UUID {
	return t.tenantUUID
}

// the alias name does not need to be unique and is alias used for quick reference when searching. no starting spaces
func (t *tenant) setAlias(alias string) error {
	// no spaces, special characters, or swear words
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(false).WithSanitizeSpecialCharacters(true).WithSanitizeAccents(false)
	if profanityDetector.IsProfane(alias) {
		err := fmt.Errorf("this is not a valid alias name, only a-z A-Z 1-9 - and spaces allowed")
		return err
	}
	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)
	if !r.MatchString(alias) {
		err := fmt.Errorf(`this is not a valid alias name `)
		return err
	}
	t.alias = alias
	return nil

}

// returns the tenant alias
func (t *tenant) getAlias() string {
	return t.alias
}

// the subdomain name must be unique within the same lord tenant. A lord tenant ties to a second level domain (eg. subscripify.com)
func (t *tenant) setSubdomainName(subdomain string) error {
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(false).WithSanitizeSpecialCharacters(true).WithSanitizeAccents(false)
	if profanityDetector.IsProfane(subdomain) {
		err := fmt.Errorf("this is not a valid subdomain name")
		return err
	}
	pattern := "^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$"
	r := regexp.MustCompile(pattern)
	if r.MatchString(subdomain) {
		t.subdomain = subdomain
		t.kubeNamespacePrefix = subdomain
		return nil
	}
	err := fmt.Errorf("not a valid subdomain name - only a-z A-Z 1-9 and - allowed")
	return err
}

// returns the tenant subdomain
func (t *tenant) getSubdomainName() string {
	return t.subdomain
}

// returns the tenant kube_namespace prefix
func (t *tenant) getKubeNamespacePrefix() string {
	return t.kubeNamespacePrefix
}

// sets the secondary domain - validates against proper naming conventions for domain names
func (t *tenant) setSecondaryDomainName(secondaryDomain string) error {

	if !(t.lordTenant) {
		err := fmt.Errorf("this value is only settable for a lord tenant - set tenant type to lord or set lord tenant value")
		return err
	} else {

		pattern := "^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$"
		r := regexp.MustCompile(pattern)
		if r.MatchString(secondaryDomain) {
			t.secondaryDomain = secondaryDomain
			return nil
		} else {
			err := fmt.Errorf("not a valid domain name - only a-z A-Z 1-9 and - allowed")
			return err
		}
	}
}

// returns the tenant secondary domain (e.g. subscripify.com - subscripify is the secondary domain)
func (t *tenant) getSecondaryDomainName() string {
	return t.secondaryDomain
}

// sets the top level domain - validates against ICANN/IANA list
func (t *tenant) setTopLevelDomain(topLevelDomain string) error {
	if !(t.lordTenant) {
		err := fmt.Errorf("this value is only settable for a lord tenant - set tenant type to lord or set lord tenant value")
		return err
	} else {

		if tld.IsValid(topLevelDomain) {
			t.topLevelDomain = topLevelDomain
			return nil
		}
		err := fmt.Errorf("not a valid top level domain - must be on this list https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
		return err
	}
}

// returns the tenant top level domain
func (t *tenant) getTopLevelDomain() string {
	return t.topLevelDomain
}

// all lord tenants require this config other types of tenants this is null - it refers to the configuration files for lord services available to this tenant
// takes in string and parses to a UUID format.
func (t *tenant) setLordServicesConfig(lordServicesConfig string) error {
	lordServiceConfigParsedUUID, err := uuid.Parse(lordServicesConfig)
	if err != nil {
		return fmt.Errorf("lord services config uuid failed to parse: %s", err)
	}

	if !t.lordTenant {
		err := fmt.Errorf("invalid tenant type for setting lord services - set tenant type to lord tenant")
		return err
	}
	t.lordServicesConfig = lordServiceConfigParsedUUID
	return nil
}

// returns the tenant lord services config
func (t *tenant) getLordServicesConfig() uuid.UUID {
	return t.lordServicesConfig
}

// all lord tenants and super tenants require this config other types of tenants this is null - it refers to the configuration files for super services available to this tenant
// takes in string and parses to a UUID format.
func (t *tenant) setSuperServicesConfig(superServicesConfig string) error {
	superServiceConfigParsedUUID, err := uuid.Parse(superServicesConfig)
	if err != nil {
		return fmt.Errorf("super services config uuid failed to parse: %s", err)
	}
	if !t.lordTenant && !t.superTenant {
		err := fmt.Errorf("invalid tenant type for setting super services - set tenant type to lord tenants or super tenant")
		return err
	}
	t.superServicesConfig = superServiceConfigParsedUUID
	return nil
}

// returns the tenant super services config
func (t *tenant) getSuperServicesConfig() uuid.UUID {
	return t.superServicesConfig
}

// all tenants require this config - it refers to the configuration files for public services available to this tenant
func (t *tenant) setPublicServicesConfig(publicServicesConfig string) error {

	publicServiceConfigParsedUUID, err := uuid.Parse(publicServicesConfig)
	if err != nil {
		return fmt.Errorf("public services config uuid failed to parse: %s", err)
	}
	t.publicServicesConfig = publicServiceConfigParsedUUID
	return nil
}

// returns the tenant public services config
func (t *tenant) getPublicServicesConfig() uuid.UUID {
	return t.publicServicesConfig
}

// all Super Tenants require this config and is set at time of creation - it refers to the private services oAuth/access configuration for the super tenant
func (t *tenant) setPrivateAccessConfig(privateAccessConfig string) error {
	privateAccessConfigParsedUUID, err := uuid.Parse(privateAccessConfig)
	if err != nil {
		return fmt.Errorf("private services config uuid failed to parse: %s", err)
	}
	if !t.superTenant {
		err := fmt.Errorf("invalid tenant type for setting private access config - set tenant type to super tenant")
		return err
	}
	t.privateAccessConfig = privateAccessConfigParsedUUID
	return nil
}

// returns the tenant private access config
func (t *tenant) getPrivateAccessConfig() uuid.UUID {
	return t.privateAccessConfig
}

// all all super and main tenants require this config and is set at time of creation - it refers to the custom services api oAuth/access configuration for this tenant
func (t *tenant) setCustomAccessConfig(customAccessConfig string) error {
	customAccessConfigParsedUUID, err := uuid.Parse(customAccessConfig)
	if err != nil {
		return fmt.Errorf("custom services config uuid failed to parse: %s", err)
	}
	if t.lordTenant {
		err := fmt.Errorf("invalid tenant type for setting custom services config - super tenants and main tenants only")
		return err
	}
	t.customAccessConfig = customAccessConfigParsedUUID
	return nil
}

// returns the tenant custom access config
func (t *tenant) getCustomAccessConfig() uuid.UUID {
	return t.customAccessConfig
}

// limited to Azure, ACS, GCP
func (t *tenant) setCloudLocation(cloudLocation CloudLocation) error {
	if cloudLocation == Azure || cloudLocation == AWS || cloudLocation == GCP {
		t.cloudLocation = cloudLocation
		return nil
	}
	err := fmt.Errorf("this is not a supported cloud provider")
	return err
}

// returns the tenant cloud location
func (t *tenant) getCloudLocation() CloudLocation {
	return t.cloudLocation
}

func (t *tenant) setLiegeUUID(liegeUUID string) error {
	liegeUUIDParsedUUID, err := uuid.Parse(liegeUUID)
	if err != nil {
		return fmt.Errorf("liege uuid failed to parse: %s", err)
	}
	if t.lordTenant {
		err := fmt.Errorf("invalid tenant type for setting liege - super tenants and main tenants only")
		return err
	}

	t.liegeUUID = liegeUUIDParsedUUID
	return nil
}

// returns the tenant liege uuid
func (t *tenant) getLiegeUUID() uuid.UUID {
	return t.liegeUUID
}

func (t *tenant) setLordUUID(lordUUID string) error {
	lordUUIDParsedUUID, err := uuid.Parse(lordUUID)
	if err != nil {
		return fmt.Errorf("custom services config uuid failed to parse: %s", err)
	}
	if t.lordTenant {
		err := fmt.Errorf("invalid tenant type for setting custom services config - super tenants and main tenants only")
		return err
	}
	t.lordUUID = lordUUIDParsedUUID
	return nil
}

// returns the tenant lord uuid
func (t *tenant) getLordUUID() uuid.UUID {
	return t.lordUUID
}

// sets the tenant type booleans in the database based upon the passed value
func (t *tenant) setTenantType(tenantType TenantType) error {

	switch tenantType {
	case LordTenant:
		t.lordTenant = true
		t.superTenant = false
	case SuperTenant:
		t.lordTenant = false
		t.superTenant = true
	case MainTenant:
		t.lordTenant = false
		t.superTenant = false
	default:
		err := fmt.Errorf("invalid tenant type")
		return err
	}
	return nil
}

// returns true if the tenant is a lord tenant
func (t *tenant) isLordTenant() bool {
	return t.lordTenant
}

// returns true if the tenant is a super tenant
func (t *tenant) isSuperTenant() bool {
	return t.superTenant
}

// id from identity system - must validate as email of current user session
func (t *tenant) setCreatedBy(userIdentifier string) error {
	t.createdBy = userIdentifier
	return nil
}

// returns the creator of the tenant
func (t *tenant) getTenantCreator() string {
	return t.createdBy
}

// returns the creation time of the tenant as recorded in the database - this value is empty for tenants yet to be inserted into the database
func (t *tenant) getCreateTime() time.Time {
	return t.createTimestamp
}
