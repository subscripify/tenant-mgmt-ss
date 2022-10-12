package tenant

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/tenantdb"
	"github.com/google/uuid"
)

type iTenantSearch interface {
	buildTenantSearchResults(rowStart int, returnCount int) (*tenantListResponseOuter, int, error)
	mapTenantUUID(pipedString string) error
	GetTenantUUIDQueryInString() string
	mapTenantAlias(pipedString string) error
	GetTenantAliasQueryLikeString() string
	mapSubdomain(pipedString string) error
	GetSubdomainQueryLikeString() string
	mapLordConfigAlias(pipedString string) error
	GetLordConfigAliasQueryLikeString() string
	mapLordConfigUUID(pipedString string) error
	GetLordConfigUUIDQueryInString() string
	mapSuperConfigAlias(pipedString string) error
	GetSuperConfigAliasQueryLikeString() string
	mapSuperConfigUUID(pipedString string) error
	GetSuperConfigUUIDQueryInString() string
	mapPublicConfigAlias(pipedString string) error
	GetPublicConfigAliasQueryLikeString() string
	mapPublicConfigUUID(pipedString string) error
	GetPublicConfigUUIDQueryInString() string
	mapPrivateAccessAlias(pipedString string) error
	GetPrivateAccessAliasQueryLikeString() string
	mapPrivateAccessUUID(pipedString string) error
	GetPrivateAccessUUIDQueryInString() string
	mapCustomAccessAlias(pipedString string) error
	GetCustomAccessAliasQueryLikeString() string
	mapCustomAccessUUID(pipedString string) error
	GetCustomAccessUUIDQueryInString() string
}

type tenantSearch struct {
	tenantType         []string
	tenantUUID         []uuid.UUID
	tenantAlias        []string
	subdomain          []string
	lordConfigAlias    []string
	lordConfigUUID     []uuid.UUID
	superConfigAlias   []string
	superConfigUUID    []uuid.UUID
	publicConfigAlias  []string
	publicConfigUUID   []uuid.UUID
	privateAccessAlias []string
	privateAccessUUID  []uuid.UUID
	customAccessAlias  []string
	customAccessUUID   []uuid.UUID
}

type tenantListResponseOuter struct {
	Results []tenantListResponseInner
}

type tenantListResponseInner struct {
	TenantUUID               string
	TenantAlias              string
	TopLevelDomain           string
	SecondaryDomain          string
	Subdomain                string
	TenantType               string
	LordConfigAlias          string
	LordServicesConfigUUID   string
	SuperConfigAlias         string
	SuperServicesConfigUUID  string
	PublicConfigAlias        string
	PublicServicesConfigUUID string
	PrivateAccessConfigAlias string
	PrivateAccessConfigUUID  string
	CustomAccessConfigAlias  string
	CustomAccessConfigUUID   string
}

func (ts *tenantSearch) buildTenantSearchResults(rowStart int, returnCount int) (*tenantListResponseOuter, int, error) {

	selectString := `SELECT BIN_TO_UUID(tenant_UUID) AS tenant_UUID, 
	tenant_alias, 
	subdomain,
	secondary_domain,
	top_level_domain, 
	tenant_type, 
	COALESCE(lord_config_alias, '') AS lord_config_alias ,
	COALESCE(BIN_TO_UUID(lord_config_UUID),'') AS lord_config_UUID, 
	COALESCE(super_config_alias, '') AS super_config_alias, 
	COALESCE(BIN_TO_UUID(super_config_UUID), '') AS super_config_UUID, 
	COALESCE(public_config_alias, '') AS public_config_alias, 
	COALESCE(BIN_TO_UUID(public_config_UUID), '') AS public_config_UUID,
	COALESCE(private_access_config_alias, '') AS private_access_config_alias, 
	COALESCE(BIN_TO_UUID(private_access_config_UUID), '') AS private_access_config_UUID, 
	COALESCE(custom_access_config_alias, '') AS custom_access_config_alias,
	COALESCE(BIN_TO_UUID(custom_access_config_UUID), '') AS custom_access_config_UUID  
	FROM tenant_search`

	isFirst := true

	t := reflect.TypeOf(ts)
	v := reflect.ValueOf(ts)

	for i := 0; i < t.NumMethod(); i++ {

		method := t.Method(i)
		if strings.HasPrefix(method.Name, "Get") {
			log.Println(method.Name)
			whereVal := v.MethodByName(method.Name).Call(nil)
			whereString := whereVal[0].String()
			if whereString != "" {
				if isFirst {
					selectString = selectString + ` WHERE `
				} else {
					selectString = selectString + ` AND `
				}
				selectString = selectString + string(whereString)
				isFirst = false
			}
		}

	}
	selectString = selectString + ` ORDER BY tenant_alias LIMIT ` + fmt.Sprint(rowStart) + `,` + fmt.Sprint(returnCount)
	log.Println(selectString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tdb := tenantdbserv.Tdb.Handle
	rows, err := tdb.QueryContext(ctx, selectString)
	if err != nil {
		return nil, 500, fmt.Errorf("invalid query: %s", err)
	}

	var sr tenantListResponseOuter
	for rows.Next() {
		var r tenantListResponseInner
		err = rows.Scan(
			&r.TenantUUID,
			&r.TenantAlias,
			&r.Subdomain,
			&r.SecondaryDomain,
			&r.TopLevelDomain,
			&r.TenantType,
			&r.LordConfigAlias,
			&r.LordServicesConfigUUID,
			&r.SuperConfigAlias,
			&r.SuperServicesConfigUUID,
			&r.PublicConfigAlias,
			&r.PublicServicesConfigUUID,
			&r.PrivateAccessConfigAlias,
			&r.PrivateAccessConfigUUID,
			&r.CustomAccessConfigAlias,
			&r.CustomAccessConfigUUID)
		if err != nil {
			return nil, 500, fmt.Errorf("invalid scan: %s", err)
		}
		sr.Results = append(sr.Results, r)
	}
	return &sr, 200, nil
}

func (ts *tenantSearch) mapTenantType(pipedString string) error {
	splitString := strings.Split(pipedString, "|")
	log.Println(splitString[0])

	pattern := `^(lord|super|main)$`
	r := regexp.MustCompile(pattern)

	for i := 0; i < len(splitString); i++ {

		if !r.MatchString(splitString[i]) {
			ts.tenantType = []string{}
			return fmt.Errorf(`this is not a valid tenant type`)
		}
		ts.tenantType = append(ts.tenantType, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetTenantTypeQueryLikeString() string {
	inString := ""
	if ts.tenantType != nil {
		inString = inString + `(tenant_type in (`
		for i, s := range ts.tenantType {
			if i == (len(ts.tenantType) - 1) {
				inString = inString + `'` + s + `'`
			} else {
				inString = inString + `'` + s + `',`
			}
		}

		inString = inString + `))`
	}
	return inString
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

func (ts *tenantSearch) GetTenantUUIDQueryInString() string {
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

	for i := 0; i < len(splitString); i++ {
		ts.tenantAlias = append(ts.tenantAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetTenantAliasQueryLikeString() string {
	inString := ""
	if ts.tenantAlias != nil {
		inString = inString + `(tenant_alias LIKE `
		for i, s := range ts.tenantAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR tenant_alias LIKE '%` + s + `%' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapSubdomain(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		ts.subdomain = append(ts.subdomain, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetSubdomainQueryLikeString() string {
	inString := ""
	if ts.subdomain != nil {
		inString = inString + `(subdomain LIKE `
		for i, s := range ts.subdomain {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR subdomain LIKE '%` + s + `%' `
			}
		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapLordConfigAlias(pipedString string) error {
	splitString := strings.Split(pipedString, "|")

	for i := 0; i < len(splitString); i++ {
		ts.lordConfigAlias = append(ts.lordConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetLordConfigAliasQueryLikeString() string {
	inString := ""
	if ts.lordConfigAlias != nil {
		inString = inString + `(lord_config_alias LIKE `
		for i, s := range ts.lordConfigAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR lord_config_alias LIKE '%` + s + `%' `
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

func (ts *tenantSearch) GetLordConfigUUIDQueryInString() string {
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

	for i := 0; i < len(splitString); i++ {
		ts.superConfigAlias = append(ts.superConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetSuperConfigAliasQueryLikeString() string {
	inString := ""
	if ts.superConfigAlias != nil {
		inString = inString + `(super_config_alias LIKE `
		for i, s := range ts.superConfigAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR super_config_alias LIKE '%` + s + `%' `
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

func (ts *tenantSearch) GetSuperConfigUUIDQueryInString() string {
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

	for i := 0; i < len(splitString); i++ {
		ts.publicConfigAlias = append(ts.publicConfigAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetPublicConfigAliasQueryLikeString() string {
	inString := ""
	if ts.publicConfigAlias != nil {
		inString = inString + `(public_config_alias LIKE `
		for i, s := range ts.publicConfigAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR public_config_alias LIKE '%` + s + `%' `
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

func (ts *tenantSearch) GetPublicConfigUUIDQueryInString() string {
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

	for i := 0; i < len(splitString); i++ {
		ts.privateAccessAlias = append(ts.privateAccessAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetPrivateAccessAliasQueryLikeString() string {
	inString := ""
	if ts.privateAccessAlias != nil {
		inString = inString + `(private_access_config_alias LIKE `
		for i, s := range ts.privateAccessAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR private_access_config_alias LIKE '%` + s + `%' `
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

func (ts *tenantSearch) GetPrivateAccessUUIDQueryInString() string {
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

	for i := 0; i < len(splitString); i++ {
		ts.customAccessAlias = append(ts.customAccessAlias, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetCustomAccessAliasQueryLikeString() string {
	inString := ""
	if ts.customAccessAlias != nil {
		inString = inString + `(custom_access_config_alias LIKE `
		for i, s := range ts.customAccessAlias {
			if i == 0 {
				inString = inString + `'%` + s + `%' `
			} else {
				inString = inString + `OR custom_access_config_alias LIKE '%` + s + `%' `
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

func (ts *tenantSearch) GetCustomAccessUUIDQueryInString() string {
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
	tenantTypes string,
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

	if tenantTypes != "" {
		err := tso.mapTenantType(tenantTypes)
		if err != nil {
			return nil, 400, err
		}
	}

	if tenantUUIDs != "" {
		err := tso.mapTenantUUID(tenantUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if tenantAliases != "" {
		err := tso.mapTenantAlias(tenantAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	if subdomains != "" {
		err := tso.mapSubdomain(subdomains)
		if err != nil {
			return nil, 400, err
		}
	}
	if lordConfigUUIDs != "" {
		err := tso.mapLordConfigUUID(lordConfigUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if lordConfigAliases != "" {
		err := tso.mapLordConfigAlias(lordConfigAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	if superConfigUUIDs != "" {
		err := tso.mapSuperConfigUUID(superConfigUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if superConfigAliases != "" {
		err := tso.mapSuperConfigAlias(superConfigAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	if publicConfigUUIDs != "" {
		err := tso.mapPublicConfigUUID(publicConfigUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if publicConfigAliases != "" {
		err := tso.mapPublicConfigAlias(publicConfigAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	if privateAccessUUIDs != "" {
		err := tso.mapPrivateAccessUUID(privateAccessUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if privateAccessAliases != "" {
		err := tso.mapPrivateAccessAlias(privateAccessAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	if customAccessUUIDs != "" {
		err := tso.mapCustomAccessUUID(customAccessUUIDs)
		if err != nil {
			return nil, 400, err
		}
	}
	if customAccessAliases != "" {
		err := tso.mapCustomAccessAlias(customAccessAliases)
		if err != nil {
			return nil, 400, err
		}
	}
	return &tso, 200, nil
}
