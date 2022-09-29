package tenant

import (
	"strings"
	"time"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/subscripify-logger"
	"github.com/google/uuid"
)

type iHttpResponse interface {
	GetHttpResponse() *httpResponseData
	generateHttpResponseCodeAndMessage(responseCode int, message string)
	generateNewTenantResponse(nID uuid.UUID)
	generateLoadedTenantResponse(l *loadedTenant, expectedAction RestAction)
}

type newTenantResponse struct {
	TenantUUID string
}

type fullTenantResponse struct {
	TenantType                         string
	TenantUUID                         uuid.UUID
	Alias                              string
	TopLevelDomain                     string
	SecondaryDomain                    string
	Subdomain                          string
	KubeNamespacePrefix                string
	SubscripifyDeploymentCloudLocation string
	LordServicesConfig                 uuid.UUID
	SuperServicesConfig                uuid.UUID
	PublicServicesConfig               uuid.UUID
	PrivateAccessConfig                uuid.UUID
	CustomAccessConfig                 uuid.UUID
	LiegeUUID                          uuid.UUID
	LordUUID                           uuid.UUID
	CreateTimestamp                    time.Time
	CreatedBy                          string
}

// type mainTenantSend struct {
// 	TenantType                         string
// 	TenantUUID                         string
// 	Alias                              string
// 	TopLevelDomain                     string
// 	SecondaryDomain                    string
// 	Subdomain                          string
// 	KubeNamespacePrefix                string
// 	SubscripifyDeploymentCloudLocation string
// 	PublicServicesConfig               string
// 	CustomAccessConfig                 string
// 	LiegeUUID                          string
// 	LordUUID                           string
// 	CreateTimestamp                    string
// 	CreatedBy                          string
// }

// type superTenantSend struct {
// 	TenantType                         string
// 	TenantUUID                         string
// 	Alias                              string
// 	TopLevelDomain                     string
// 	SecondaryDomain                    string
// 	Subdomain                          string
// 	KubeNamespacePrefix                string
// 	SubscripifyDeploymentCloudLocation string
// 	SuperServicesConfig                string
// 	PublicServicesConfig               string
// 	PrivateAccessConfig                string
// 	CustomAccessConfig                 string
// 	LiegeUUID                          string
// 	LordUUID                           string
// 	CreateTimestamp                    string
// 	CreatedBy                          string
// }

// type lordTenantSend struct {
// 	TenantType                         string
// 	TenantUUID                         string
// 	Alias                              string
// 	TopLevelDomain                     string
// 	SecondaryDomain                    string
// 	Subdomain                          string
// 	KubeNamespacePrefix                string
// 	SubscripifyDeploymentCloudLocation string
// 	LordServiceConfig                  string
// 	SuperServicesConfig                string
// 	PublicServicesConfig               string
// 	CreateTimestamp                    string
// 	CreatedBy                          string
// }

type httpResponseData struct {
	HttpResponseCode int
	Message          string
	NewTenant        newTenantResponse
	FullTenant       fullTenantResponse
}

type RestAction string

const (
	GET    RestAction = "GET"
	POST   RestAction = "POST"
	PUT    RestAction = "PUT"
	PATCH  RestAction = "PATCH"
	DELETE RestAction = "DELETE"
)

func (r *httpResponseData) GetHttpResponse() *httpResponseData {
	return r
}

// this function generates the data for http responses, regardless if the api uses them in a response or not
func (hr *httpResponseData) generateHttpResponseCodeAndMessage(responseCode int, message string) {

	hr.HttpResponseCode = responseCode
	hr.Message = message

	if responseCode >= 400 && responseCode <= 499 {
		subscripifylogger.DebugLog.Printf(
			"tenant service generated http response : %v %s ",
			responseCode,
			strings.ToLower(message))
		return
	}
	if responseCode >= 500 && responseCode <= 599 {
		subscripifylogger.ErrorLog.Printf(
			"tenant service generated http response : %v %s",
			responseCode,
			strings.ToLower(message))
		return
	}
	subscripifylogger.DebugLog.Printf(
		"tenant service generated http response : %v %s ",
		responseCode,
		strings.ToLower(message))

}

// this function generates a new tenant response and assigns it to a
func (tr *httpResponseData) generateNewTenantResponse(nID uuid.UUID) {

	tr.NewTenant.TenantUUID = nID.String()

	subscripifylogger.DebugLog.Printf(
		"tenant service sent new tenant UUID:%s ",
		nID.String())

}

func (tr *httpResponseData) generateLoadedTenantResponse(l *loadedTenant, expectedAction RestAction) {
	if l.isLordTenant() {
		tr.FullTenant.TenantType = string(LordTenant)
	} else if l.isSuperTenant() {
		tr.FullTenant.TenantType = string(SuperTenant)
	} else {
		tr.FullTenant.TenantType = string(MainTenant)
	}

	tr.FullTenant.TenantUUID = l.getTenantUUID()
	tr.FullTenant.Alias = l.getAlias()
	tr.FullTenant.TopLevelDomain = l.getTopLevelDomain()
	tr.FullTenant.SecondaryDomain = l.getSecondaryDomainName()
	tr.FullTenant.Subdomain = l.getSubdomainName()
	tr.FullTenant.KubeNamespacePrefix = l.getKubeNamespacePrefix()
	tr.FullTenant.SubscripifyDeploymentCloudLocation = string(l.getCloudLocation())
	tr.FullTenant.LordServicesConfig = l.getLordServicesConfig()
	tr.FullTenant.SuperServicesConfig = l.getSuperServicesConfig()
	tr.FullTenant.PublicServicesConfig = l.getPublicServicesConfig()
	tr.FullTenant.PrivateAccessConfig = l.getPrivateAccessConfig()
	tr.FullTenant.CustomAccessConfig = l.getCustomAccessConfig()
	tr.FullTenant.CreateTimestamp = l.getCreateTime()
	tr.FullTenant.CreatedBy = l.getTenantCreator()
	tr.FullTenant.LiegeUUID = l.getLiegeUUID()
	tr.FullTenant.LordUUID = l.getLordUUID()

	subscripifylogger.DebugLog.Printf(
		"tenant service sent a full tenant object for a  %s action :%s ",
		expectedAction,
		l.getTenantUUID().String())

}
