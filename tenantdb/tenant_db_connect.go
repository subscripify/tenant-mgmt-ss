package tenantdbserv

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"os"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/_git/tenant-mgmt-ss/subscripify-logger"
	"github.com/go-sql-driver/mysql"
)

type TenantDb interface {
	PingDb() error
}

type tenantDb struct {
	Handle *sql.DB
}

var Tdb tenantDb

func init() {
	Tdb = tenantDb{Handle: getNewMySQLTenantDbHandle()}
}

// builds a connection to an SQL database
func getNewMySQLTenantDbHandle() *sql.DB {
	var tenantsDbHandle *sql.DB

	rootCertPool := x509.NewCertPool()
	pem, _ := ioutil.ReadFile(os.Getenv("DBAPPCERTLOCATION")) // this is the only one that one can use with AzureMYSQL

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		subscripifylogger.FatalLog.Fatal("Failed to append PEM.")
	}
	env := os.Getenv("SUBSCRIPIFY_DB_ENV")
	var cfgdb mysql.Config
	if !(env == "localdb") {

		mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})

		cfgdb = mysql.Config{
			User:                 os.Getenv("DBUSER"),
			Passwd:               os.Getenv("DBPASS"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DBHOST"),
			DBName:               "tenants",
			AllowNativePasswords: true,
			TLSConfig:            "custom",
			ParseTime:            true,
		}
	} else {

		cfgdb = mysql.Config{
			User:   "root",
			Passwd: "insecure",
			Net:    "tcp",
			Addr:   "localhost",

			DBName:               "tenants",
			AllowNativePasswords: true,
			ParseTime:            true,
		}
	}
	// Get a database handle.
	var err error
	tenantsDbHandle, err = sql.Open("mysql", cfgdb.FormatDSN())

	if err != nil {
		subscripifylogger.FatalLog.Fatal(err)
	}

	pingErr := tenantsDbHandle.Ping()
	if pingErr != nil {
		subscripifylogger.FatalLog.Fatal(pingErr)
	}
	if !(env == "localdb") {
		subscripifylogger.InfoLog.Printf("connected to tenants db on server:%s", "localhost")
	} else {
		subscripifylogger.InfoLog.Printf("connected to tenants db on server:%s", os.Getenv("DBHOST"))
	}

	return tenantsDbHandle
}

func (tdb *tenantDb) PingDb() error {
	err := tdb.Handle.Ping()
	if err != nil {
		subscripifylogger.FatalLog.Printf(err.Error())
		return err
	}
	return nil
}
