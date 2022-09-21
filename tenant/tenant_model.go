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
	setCreateTimestamp()
}

type tenant struct {
	tenantUUID           uuid.UUID
	alias                string
	subdomain            string
	secondaryDomain      string
	topLevelDomain       string
	kubeNamespacePrefix  string
	lordServicesConfig   uuid.UUID
	superServicesConfig  uuid.UUID
	publicServicesConfig uuid.UUID
	privateAccessConfig  uuid.UUID
	customAccessConfig   uuid.UUID
	cloudLocation        CloudLocation
	liegeUUID            uuid.UUID
	lordUUID             uuid.UUID
	isLordTenant         bool
	isSuperTenant        bool
	createTimestamp      time.Time
	createdBy            string
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
	ACS   CloudLocation = "ACS"
	GCP   CloudLocation = "GCP"
)

// this is the only place in the application this value is created
func (t *tenant) setNewTenantUUID() {
	t.tenantUUID = uuid.New()
}

// the org name does not need to be unique and is alias used for quick reference when searching. no starting spaces
func (t *tenant) setAlias(alias string) error {
	// no spaces, special characters, or swear words
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(false).WithSanitizeSpecialCharacters(true).WithSanitizeAccents(false)
	if profanityDetector.IsProfane(alias) {
		err := fmt.Errorf("this is not a valid alias name")
		return err
	}
	pattern := `(?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)
	if r.MatchString(alias) {
		t.alias = alias
		return nil
	}
	err := fmt.Errorf(`this is not a valid alias name must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
	return err

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
	err := fmt.Errorf("not a valid subdomain name - must match pattern '^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$'")
	return err
}

// sets the secondary domain - validates against proper naming conventions for domain names
func (t *tenant) setSecondaryDomainName(secondaryDomain string) error {

	if !(t.isLordTenant) {
		err := fmt.Errorf("this value is only settable for a lord tenant - set tenant type to lord or set lord tenant value")
		return err
	} else {

		pattern := "^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$"
		r := regexp.MustCompile(pattern)
		if r.MatchString(secondaryDomain) {
			t.secondaryDomain = secondaryDomain
			return nil
		} else {
			err := fmt.Errorf("not a valid domain name - must match pattern '^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$'")
			return err
		}
	}
}

// sets the top level domain - validates against ICANN/IANA list
func (t *tenant) setTopLevelDomain(topLevelDomain string) error {
	if !(t.isLordTenant) {
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

// all lord tenants require this config other types of tenants this is null - it refers to the configuration files for lord services available to this tenant
// takes in string and parses to a UUID format.
func (t *tenant) setLordServicesConfig(lordServicesConfig string) error {
	lordServiceConfigParsedUUID, err := uuid.Parse(lordServicesConfig)
	if err != nil {
		return fmt.Errorf("lord services config uuid failed to parse: %s", err)
	}

	if !t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting lord services - set tenant type to lord tenant")
		return err
	}
	t.lordServicesConfig = lordServiceConfigParsedUUID
	return nil
}

// all lord tenants and super tenants require this config other types of tenants this is null - it refers to the configuration files for super services available to this tenant
// takes in string and parses to a UUID format.
func (t *tenant) setSuperServicesConfig(superServicesConfig string) error {
	superServiceConfigParsedUUID, err := uuid.Parse(superServicesConfig)
	if err != nil {
		return fmt.Errorf("super services config uuid failed to parse: %s", err)
	}
	if !t.isLordTenant && !t.isSuperTenant {
		err := fmt.Errorf("invalid tenant type for setting super services - set tenant type to lord tenants or super tenant")
		return err
	}
	t.superServicesConfig = superServiceConfigParsedUUID
	return nil
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

// all all Super Tenants require this config and is set at time of creation - it refers to the private services oAuth/access configuration for the super tenant
func (t *tenant) setPrivateAccessConfig(privateAccessConfig string) error {
	privateAccessConfigParsedUUID, err := uuid.Parse(privateAccessConfig)
	if err != nil {
		return fmt.Errorf("private services config uuid failed to parse: %s", err)
	}
	if !t.isSuperTenant {
		err := fmt.Errorf("invalid tenant type for setting private access config - set tenant type to super tenant")
		return err
	}
	t.privateAccessConfig = privateAccessConfigParsedUUID
	return nil
}

// all all super and main tenants require this config and is set at time of creation - it refers to the custom services api oAuth/access configuration for this tenant
func (t *tenant) setCustomAccessConfig(customAccessConfig string) error {
	customAccessConfigParsedUUID, err := uuid.Parse(customAccessConfig)
	if err != nil {
		return fmt.Errorf("custom services config uuid failed to parse: %s", err)
	}
	if t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting custom services config - super tenants and main tenants only")
		return err
	}
	t.customAccessConfig = customAccessConfigParsedUUID
	return nil
}

// limited to Azure, ACS, GCP
func (t *tenant) setCloudLocation(cloudLocation CloudLocation) error {
	if cloudLocation == Azure || cloudLocation == ACS || cloudLocation == GCP {
		t.cloudLocation = cloudLocation
		return nil
	}
	err := fmt.Errorf("this is not a supported cloud provider")
	return err
}

func (t *tenant) setLiegeUUID(liegeUUID string) error {
	liegeUUIDParsedUUID, err := uuid.Parse(liegeUUID)
	if err != nil {
		return fmt.Errorf("custom services config uuid failed to parse: %s", err)
	}
	if t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting custom services config - super tenants and main tenants only")
		return err
	}
	t.liegeUUID = liegeUUIDParsedUUID
	return nil
}

func (t *tenant) setLordUUID(lordUUID string) error {
	lordUUIDParsedUUID, err := uuid.Parse(lordUUID)
	if err != nil {
		return fmt.Errorf("custom services config uuid failed to parse: %s", err)
	}
	if t.isLordTenant {
		err := fmt.Errorf("invalid tenant type for setting custom services config - super tenants and main tenants only")
		return err
	}
	t.lordUUID = lordUUIDParsedUUID
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
		err := fmt.Errorf("invalid tenant type")
		return err
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
