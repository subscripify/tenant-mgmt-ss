package tenant

import (
	"context"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

	tenantdbserv "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/tenantdb"
	"github.com/google/uuid"
)

type iTenantSearch interface {
	buildTenantSearchResults(rowStart int, returnCount int) (*tenantListResponseOuter, int, error)
	mapTenantUUID() error
	GetTenantUUIDQueryInString() string
	mapTenantAlias() error
	GetTenantAliasQueryLikeString() string
	mapSubdomain() error
	GetSubdomainQueryLikeString() string
	mapDomain() error
	GetDomainQueryLikeString() string
	mapLordConfigAlias() error
	GetLordConfigAliasQueryLikeString() string
	mapLordConfigUUID() error
	GetLordConfigUUIDQueryInString() string
	mapSuperConfigAlias() error
	GetSuperConfigAliasQueryLikeString() string
	mapSuperConfigUUID() error
	GetSuperConfigUUIDQueryInString() string
	mapPublicConfigAlias() error
	GetPublicConfigAliasQueryLikeString() string
	mapPublicConfigUUID() error
	GetPublicConfigUUIDQueryInString() string
	mapPrivateAccessAlias() error
	GetPrivateAccessAliasQueryLikeString() string
	mapPrivateAccessUUID() error
	GetPrivateAccessUUIDQueryInString() string
	mapCustomAccessAlias() error
	GetCustomAccessAliasQueryLikeString() string
	mapCustomAccessUUID() error
	GetCustomAccessUUIDQueryInString() string
}

type tenantSearch struct {
	tenantType         []string
	tenantUUID         []uuid.UUID
	tenantAlias        []string
	subdomain          []string
	domain             []string
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
	qsType             string
	qsTid              string
	qsTal              string
	qsSubdmn           string
	qsDmn              string
	qsCal              string
	qsCid              string
	qsAal              string
	qsAid              string
}

type tenantListResponseOuter struct {
	Paging *tenantSearchResultsPaging
	Data   []tenantListResponseInner
}

type tenantSearchResultsPaging struct {
	PageCount int
	RowCount  int
	Previous  string
	Next      string
}

type tenantListResponseInner struct {
	TenantUUID         string
	TenantAlias        string
	TenantSubdomain    string
	TenantSecDomain    string
	TenantTld          string
	TenantType         string
	LordConfigAlias    string
	LordConfigUUID     string
	SuperConfigAlias   string
	SuperConfigUUID    string
	PublicConfigAlias  string
	PublicConfigUUID   string
	PrivateAccessAlias string
	PrivateAccessUUID  string
	CustomAccessAlias  string
	CustomAccessUUID   string
}

func (ts *tenantSearch) buildTenantSearchResults(page int, perPage int) (*tenantListResponseOuter, int, error) {

	if page < 1 {
		return nil, 400, fmt.Errorf("page must be an integer greater than 1")
	}
	if perPage <= 0 {
		return nil, 400, fmt.Errorf("per page must be an integer greater than 0")
	}

	selectString := `SELECT 
	COUNT(tenant_uuid) OVER() as total_count,
	BIN_TO_UUID(tenant_UUID) AS tenant_UUID, 
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
	selectString = selectString + ` ORDER BY tenant_alias LIMIT ` + fmt.Sprint((page-1)*perPage) + `,` + fmt.Sprint(perPage)
	log.Println(selectString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tdb := tenantdbserv.Tdb.Handle
	rows, err := tdb.QueryContext(ctx, selectString)

	if err != nil {
		return nil, 500, fmt.Errorf("invalid query: %s", err)
	}

	var itemCount int32
	var sr tenantListResponseOuter

	for rows.Next() {
		var r tenantListResponseInner
		err = rows.Scan(
			&itemCount,
			&r.TenantUUID,
			&r.TenantAlias,
			&r.TenantSubdomain,
			&r.TenantSecDomain,
			&r.TenantTld,
			&r.TenantType,
			&r.LordConfigAlias,
			&r.LordConfigUUID,
			&r.SuperConfigAlias,
			&r.SuperConfigUUID,
			&r.PublicConfigAlias,
			&r.PublicConfigUUID,
			&r.PrivateAccessAlias,
			&r.PrivateAccessUUID,
			&r.CustomAccessAlias,
			&r.CustomAccessUUID)

		sr.Data = append(sr.Data, r)
	}
	if err != nil {
		return nil, 500, fmt.Errorf("invalid scan: %s", err)
	}
	if len(sr.Data) == 0 {

		if page > 1 {
			return nil, 404, fmt.Errorf("results are zero and requested page is greater than 1. re-try with pg=1 in query parameters to verify search criteria")
		}
	}
	sr.Paging = ts.buildPagination(int(itemCount), page, perPage)
	return &sr, 200, nil
}

func (ts *tenantSearch) buildPagination(rowCount int, page int, perPage int) *tenantSearchResultsPaging {
	var paginationObject tenantSearchResultsPaging
	if rowCount == 0 {

		return &paginationObject
	}
	if perPage > 0 {
		divide := (float64(rowCount) / float64(perPage))
		paginationObject.PageCount = int(math.Ceil(divide))
	} else {
		paginationObject.PageCount = 1
	}
	var nextPage int
	var lastPage bool
	if page >= paginationObject.PageCount {
		nextPage = paginationObject.PageCount
		lastPage = true
	} else {
		nextPage = page + 1
	}
	var prevPage int
	var firstPage bool
	if page <= 1 {
		prevPage = 1
		firstPage = true
	} else {
		prevPage = page - 1
	}

	var queryString string
	if ts.qsType != "" {
		queryString = `&type=` + ts.qsType
	}
	if ts.qsTid != "" {
		queryString = `&tid=` + ts.qsTid
	}
	if ts.qsTal != "" {
		queryString = `&tal=` + ts.qsTal
	}
	if ts.qsSubdmn != "" {
		queryString = `&subdmn=` + ts.qsSubdmn
	}
	if ts.qsDmn != "" {
		queryString = `&dmn=` + ts.qsDmn
	}
	if ts.qsCal != "" {
		queryString = `&cal=` + ts.qsCal
	}
	if ts.qsCid != "" {
		queryString = `&cid=` + ts.qsCid
	}
	if ts.qsAal != "" {
		queryString = `&aal=` + ts.qsAal
	}
	if ts.qsAid != "" {
		queryString = `&aid=` + ts.qsAid
	}

	var nextPageQueryString string
	var prevPageQueryString string
	if !lastPage {
		nextPageQueryString = `/search/tenants?pg=` + fmt.Sprint(nextPage)
		nextPageQueryString = nextPageQueryString + `&lc=` + fmt.Sprint(perPage)
		nextPageQueryString = nextPageQueryString + queryString
	}
	if !firstPage {
		prevPageQueryString = `/search/tenants?pg=` + fmt.Sprint(prevPage)
		prevPageQueryString = prevPageQueryString + `&lc=` + fmt.Sprint(perPage)
		prevPageQueryString = prevPageQueryString + queryString
	}
	paginationObject.RowCount = rowCount
	paginationObject.Next = nextPageQueryString
	paginationObject.Previous = prevPageQueryString

	return &paginationObject
}

func (ts *tenantSearch) mapTenantType() error {
	splitString := strings.Split(ts.qsType, "|")
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

func (ts *tenantSearch) mapTenantUUID() error {
	splitString := strings.Split(ts.qsTid, "|")

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

func (ts *tenantSearch) mapTenantAlias() error {
	splitString := strings.Split(ts.qsTal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapSubdomain() error {
	splitString := strings.Split(ts.qsSubdmn, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for subdomain only a-z A-Z 1-9 and - allowed")
		}
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

func (ts *tenantSearch) mapDomain() error {
	splitString := strings.Split(ts.qsDmn, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for domain only a-z A-Z 1-9 and - allowed")
		}
		ts.domain = append(ts.domain, splitString[i])
	}

	return nil
}

func (ts *tenantSearch) GetDomainQueryLikeString() string {
	inString := ""
	topDmn := ""
	secDmn := ""
	log.Println(ts.domain)
	if ts.domain != nil {
		for i, s := range ts.domain {
			secTopSplit := strings.Split(s, ".")
			log.Println(secTopSplit)
			if len(secTopSplit) < 2 {
				topDmn = ""
				secDmn = secTopSplit[0]
			} else {
				topDmn = secTopSplit[len(secTopSplit)-1]
				secDmn = strings.Join(secTopSplit[0:(len(secTopSplit)-1)], ".")
			}
			if topDmn == "" {
				if i == 0 {
					inString = inString + `(secondary_domain LIKE '%` + secDmn + `%' `
				} else {
					inString = inString + `OR secondary_domain LIKE '%` + secDmn + `%' `
				}
			} else {
				if i == 0 {
					inString = inString + `((secondary_domain LIKE '%` + secDmn + `%' AND top_level_domain LIKE '%` + topDmn + `%') `
				} else {
					inString = inString + `OR (secondary_domain LIKE '%` + secDmn + `%' AND top_level_domain LIKE '%` + topDmn + `%') `
				}
			}

		}
		inString = inString + `)`
	}
	return inString
}

func (ts *tenantSearch) mapLordConfigAlias() error {
	splitString := strings.Split(ts.qsCal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapLordConfigUUID() error {
	splitString := strings.Split(ts.qsCid, "|")

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

func (ts *tenantSearch) mapSuperConfigAlias() error {
	splitString := strings.Split(ts.qsCal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapSuperConfigUUID() error {
	splitString := strings.Split(ts.qsCid, "|")

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

func (ts *tenantSearch) mapPublicConfigAlias() error {
	splitString := strings.Split(ts.qsCal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapPublicConfigUUID() error {
	splitString := strings.Split(ts.qsCid, "|")

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

func (ts *tenantSearch) mapPrivateAccessAlias() error {
	splitString := strings.Split(ts.qsAal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapPrivateAccessUUID() error {
	splitString := strings.Split(ts.qsAid, "|")

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

func (ts *tenantSearch) mapCustomAccessAlias() error {
	splitString := strings.Split(ts.qsAal, "|")

	for i := 0; i < len(splitString); i++ {
		pattern := `^[a-zA-Z0-9\s-]*$`
		re := regexp.MustCompile(pattern)
		if !re.Match([]byte(splitString[i])) {
			return fmt.Errorf("invalid search string for alias only a-z A-Z 1-9 - and spaces allowed")
		}
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

func (ts *tenantSearch) mapCustomAccessUUID() error {
	splitString := strings.Split(ts.qsAid, "|")

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
	domains string,
	configUUIDs string,
	configAliases string,
	accessUUIDs string,
	accessAliases string) (iTenantSearch, int, error) {
	var tso tenantSearch

	tso.qsType = tenantTypes
	tso.qsTid = tenantUUIDs
	tso.qsTal = tenantAliases
	tso.qsSubdmn = subdomains
	tso.qsDmn = domains
	tso.qsCal = configAliases
	tso.qsCid = configUUIDs
	tso.qsAal = accessAliases
	tso.qsAid = accessUUIDs

	if tenantTypes != "" {
		err := tso.mapTenantType()
		if err != nil {
			return nil, 400, err
		}
	}

	if tenantUUIDs != "" {
		err := tso.mapTenantUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if tenantAliases != "" {
		err := tso.mapTenantAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	if subdomains != "" {
		err := tso.mapSubdomain()
		if err != nil {
			return nil, 400, err
		}
	}
	if domains != "" {
		err := tso.mapDomain()
		if err != nil {
			return nil, 400, err
		}
	}
	if configUUIDs != "" {
		err := tso.mapLordConfigUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if configAliases != "" {
		err := tso.mapLordConfigAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	if configUUIDs != "" {
		err := tso.mapSuperConfigUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if configAliases != "" {
		err := tso.mapSuperConfigAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	if configUUIDs != "" {
		err := tso.mapPublicConfigUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if configAliases != "" {
		err := tso.mapPublicConfigAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	if accessUUIDs != "" {
		err := tso.mapPrivateAccessUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if accessAliases != "" {
		err := tso.mapPrivateAccessAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	if accessUUIDs != "" {
		err := tso.mapCustomAccessUUID()
		if err != nil {
			return nil, 400, err
		}
	}
	if accessAliases != "" {
		err := tso.mapCustomAccessAlias()
		if err != nil {
			return nil, 400, err
		}
	}
	return &tso, 200, nil
}
