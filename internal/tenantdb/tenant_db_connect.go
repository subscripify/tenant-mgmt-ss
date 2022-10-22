package tenantdbserv

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"log"
	"os"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/subscripify-logger.git"
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

	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile("publicpem/BaltimoreCyberTrustRoot.crt.pem") // this is the only one that one can use with AzureMYSQL
	if err != nil {
		subscripifylogger.FatalLog.Fatalf("Failed to read PEM: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		subscripifylogger.FatalLog.Fatal("Failed to append PEM.")
	}
	env := os.Getenv("SUBSCRIPIFY_DB_ENV")
	var cfgdb mysql.Config
	if env == "proddb" {

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
	} else if env == "localdb" {

		cfgdb = mysql.Config{
			User:   "root",
			Passwd: "insecure",
			Net:    "tcp",
			Addr:   "localhost",

			DBName:               "tenants",
			AllowNativePasswords: true,
			ParseTime:            true,
		}
	} else {
		log.Fatal("SUBSCRIPIFY_DB_ENV is not set or invalid")
	}
	// Get a database handle.

	tenantsDbHandle, err := sql.Open("mysql", cfgdb.FormatDSN())
	log.Print(cfgdb.FormatDSN())
	if err != nil {
		subscripifylogger.FatalLog.Fatal(err)
	}

	pingErr := tenantsDbHandle.Ping()
	if pingErr != nil {
		subscripifylogger.FatalLog.Fatal(pingErr)
	}
	if env == "proddb" {
		subscripifylogger.InfoLog.Printf("connected to tenants db on server:%s", os.Getenv("DBHOST"))
	} else {
		subscripifylogger.InfoLog.Printf("connected to tenants db on server:%s", "localhost")
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
