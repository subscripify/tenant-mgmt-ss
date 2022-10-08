package tenant

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type iTenantSearch interface {
	mapTenantUUID(pipedString string) error
	getTenantUUIDQueryInString() string
	mapTenantAlias(pipedString string) error
	getTenantAliasQueryLikeString() string
	mapSubdomain(pipedString string) error
	getSubdomainQueryLikeString() string
	mapLordConfigAlias(pipedString string) error
	getLordConfigAliasQueryLikeString() string
	mapLordConfigUUID(pipedString string) error
	getLordConfigUUIDQueryInString() string
	mapSuperConfigAlias(pipedString string) error
	getSuperConfigAliasQueryLikeString() string
	mapSuperConfigUUID(pipedString string) error
	getSuperConfigUUIDQueryInString() string
	mapPublicConfigAlias(pipedString string) error
	getPublicConfigAliasQueryLikeString() string
	mapPublicConfigUUID(pipedString string) error
	getPublicConfigUUIDQueryInString() string
	mapPrivateAccessAlias(pipedString string) error
	getPrivateAccessAliasQueryLikeString() string
	mapPrivateAccessUUID(pipedString string) error
	getPrivateAccessUUIDQueryInString() string
	mapCustomAccessAlias(pipedString string) error
	getCustomAccessAliasQueryLikeString() string
	mapCustomAccessUUID(pipedString string) error
	getCustomAccessUUIDQueryInString() string
}

type tenantSearch struct {
	tenantUUID        []uuid.UUID
	tenantAlias       []string
	subdomain         []string
	lordConfigAlias   []string
	lordConfigUUID    []uuid.UUID
	superConfigAlias  []string
	superConfigUUID   []uuid.UUID
	publicConfigAlias []string
	publicConfigUUID  []uuid.UUID

	privateAccessAlias []string
	privateAccessUUID  []uuid.UUID
	customAccessAlias  []string
	customAccessUUID   []uuid.UUID
}

func (ts *tenantSearch) mapTenantUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.tenantUUID = []uuid.UUID{}
			return fmt.Errorf("tenant uuid failed to parse: %s", err)
		}
		ts.tenantUUID = append(ts.tenantUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getTenantUUIDQueryInString() string {
	inString := ""
	if ts.tenantUUID != nil {
		inString = inString + `(tenant_uuid IN (`
		for i, s := range ts.tenantUUID {
			if i == (len(ts.tenantUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapTenantAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.tenantAlias = []string{}
			return fmt.Errorf(`this is not a valid tenant alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.tenantAlias = append(ts.tenantAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getTenantAliasQueryLikeString() string {
	inString := ""
	if ts.tenantAlias != nil {
		inString = inString + `(tenant_alias LIKE `
		for i, s := range ts.tenantAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR tenant_alias LIKE '` + s + `' `
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapSubdomain(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.subdomain = []string{}
			return fmt.Errorf("not a valid subdomain name - must match pattern '^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.-]*[a-zA-Z0-9]+))$'")
		}
		ts.subdomain = append(ts.subdomain, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getSubdomainQueryLikeString() string {
	inString := ""
	if ts.subdomain != nil {
		inString = inString + `subdomain LIKE `
		for i, s := range ts.subdomain {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR subdomain LIKE '` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapLordConfigAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.lordConfigAlias = []string{}
			return fmt.Errorf(`this is not a valid config alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.lordConfigAlias = append(ts.lordConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getLordConfigAliasQueryLikeString() string {
	inString := ""
	if ts.lordConfigAlias != nil {
		inString = inString + `lord_config_alias LIKE `
		for i, s := range ts.lordConfigAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR lord_config_alias LIKE ` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapLordConfigUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.lordConfigUUID = []uuid.UUID{}
			return fmt.Errorf("config uuid failed to parse: %s", err)
		}
		ts.lordConfigUUID = append(ts.lordConfigUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getLordConfigUUIDQueryInString() string {
	inString := ""
	if ts.lordConfigUUID != nil {
		inString = inString + `(lord_config_uuid IN (`
		for i, s := range ts.lordConfigUUID {
			if i == (len(ts.lordConfigUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapSuperConfigAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.superConfigAlias = []string{}
			return fmt.Errorf(`this is not a valid config alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.superConfigAlias = append(ts.superConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getSuperConfigAliasQueryLikeString() string {
	inString := ""
	if ts.superConfigAlias != nil {
		inString = inString + `super_config_alias LIKE `
		for i, s := range ts.superConfigAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR super_config_alias LIKE ` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapSuperConfigUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.superConfigUUID = []uuid.UUID{}
			return fmt.Errorf("config uuid failed to parse: %s", err)
		}
		ts.superConfigUUID = append(ts.superConfigUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getSuperConfigUUIDQueryInString() string {
	inString := ""
	if ts.superConfigUUID != nil {
		inString = inString + `(super_config_uuid IN (`
		for i, s := range ts.superConfigUUID {
			if i == (len(ts.superConfigUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapPublicConfigAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.publicConfigAlias = []string{}
			return fmt.Errorf(`this is not a valid config alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.publicConfigAlias = append(ts.publicConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getPublicConfigAliasQueryLikeString() string {
	inString := ""
	if ts.publicConfigAlias != nil {
		inString = inString + `public_config_alias LIKE `
		for i, s := range ts.publicConfigAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR public_config_alias LIKE ` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapPublicConfigUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.publicConfigUUID = []uuid.UUID{}
			return fmt.Errorf("config uuid failed to parse: %s", err)
		}
		ts.publicConfigUUID = append(ts.publicConfigUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getPublicConfigUUIDQueryInString() string {
	inString := ""
	if ts.publicConfigUUID != nil {
		inString = inString + `(public_config_uuid IN (`
		for i, s := range ts.publicConfigUUID {
			if i == (len(ts.publicConfigUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapPrivateAccessAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.privateAccessAlias = []string{}
			return fmt.Errorf(`this is not a valid access alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.privateAccessAlias = append(ts.privateAccessAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getPrivateAccessAliasQueryLikeString() string {
	inString := ""
	if ts.privateAccessAlias != nil {
		inString = inString + `private_access_alias LIKE `
		for i, s := range ts.privateAccessAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR private_access_alias LIKE ` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapPrivateAccessUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.privateAccessUUID = []uuid.UUID{}
			return fmt.Errorf("access uuid failed to parse: %s", err)
		}
		ts.privateAccessUUID = append(ts.privateAccessUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getPrivateAccessUUIDQueryInString() string {
	inString := ""
	if ts.privateAccessUUID != nil {
		inString = inString + `(private_access_config_uuid IN (`
		for i, s := range ts.privateAccessUUID {
			if i == (len(ts.privateAccessUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func (ts *tenantSearch) mapCustomAccessAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	pattern := `^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if r.MatchString(splitString[i]) {
			ts.customAccessAlias = []string{}
			return fmt.Errorf(`this is not a valid access alias must match pattern (?m)^([a-zA-Z0-9]|(?:[a-zA-Z0-9]+[a-zA-Z0-9.\s\-]*[a-zA-Z0-9]+))$`)
		}
		ts.customAccessAlias = append(ts.customAccessAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) getCustomAccessAliasQueryLikeString() string {
	inString := ""
	if ts.customAccessAlias != nil {
		inString = inString + `custom_access_alias LIKE `
		for i, s := range ts.customAccessAlias {
			if i == 0 {
				inString = inString + `'` + s + `' `
			} else {
				inString = inString + `OR custom_access_alias LIKE ` + s + `' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapCustomAccessUUID(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		parsedUUID, err := uuid.Parse(splitString[i])
		if err != nil {
			ts.customAccessUUID = []uuid.UUID{}
			return fmt.Errorf("access uuid failed to parse: %s", err)
		}
		ts.customAccessUUID = append(ts.customAccessUUID, parsedUUID)
	}

	return nil
}

func (ts *tenantSearch) getCustomAccessUUIDQueryInString() string {
	inString := ""
	if ts.customAccessUUID != nil {
		inString = inString + `(custom_access_config_uuid IN (`
		for i, s := range ts.customAccessUUID {
			if i == (len(ts.customAccessUUID) - 1) {
				inString = inString + `UUID_TO_BIN('` + s.String() + `')`
			} else {
				inString = inString + `UUID_TO_BIN('` + s.String() + `'),`
			}
		}
		inString = inString + `))`
	}
	return inString
}

func createTenantSearchObject(
	tenantUUIDs string,
	tenantAliases string,
	subdomains string,
	lordConfigUUIDs string,
	lordConfigAliases string,
	superConfigUUIDs string,
	superConfigAliases string,
	publicConfigUUIDs string,
	publicConfigAliases string,
	privateAccessUUIDs string,
	privateAccessAliases string,
	customAccessUUIDs string,
	customAccessAliases string) (iTenantSearch, int, error) {
	var tso tenantSearch

	if tenantUUIDs != "" {
		err := tso.mapTenantUUID(tenantUUIDs)
		return nil, 400, err
	}
	if tenantAliases != "" {
		err := tso.mapTenantAlias(tenantAliases)
		return nil, 400, err
	}
	if subdomains != "" {
		err := tso.mapSubdomain(subdomains)
		return nil, 400, err
	}
	if lordConfigUUIDs != "" {
		err := tso.mapLordConfigUUID(lordConfigUUIDs)
		return nil, 400, err
	}
	if lordConfigAliases != "" {
		err := tso.mapLordConfigAlias(lordConfigAliases)
		return nil, 400, err
	}
	if superConfigUUIDs != "" {
		err := tso.mapSuperConfigUUID(superConfigUUIDs)
		return nil, 400, err
	}
	if superConfigAliases != "" {
		err := tso.mapSuperConfigAlias(superConfigAliases)
		return nil, 400, err
	}
	if publicConfigUUIDs != "" {
		err := tso.mapPublicConfigUUID(publicConfigUUIDs)
		return nil, 400, err
	}
	if publicConfigAliases != "" {
		err := tso.mapPublicConfigAlias(publicConfigAliases)
		return nil, 400, err
	}
	if privateAccessUUIDs != "" {
		err := tso.mapPrivateAccessUUID(privateAccessUUIDs)
		return nil, 400, err
	}
	if privateAccessAliases != "" {
		err := tso.mapPrivateAccessAlias(privateAccessAliases)
		return nil, 400, err
	}
	if customAccessUUIDs != "" {
		err := tso.mapCustomAccessUUID(customAccessUUIDs)
		return nil, 400, err
	}
	if customAccessAliases != "" {
		err := tso.mapCustomAccessAlias(customAccessAliases)
		return nil, 400, err
	}
	return &tso, 200, nil
}
