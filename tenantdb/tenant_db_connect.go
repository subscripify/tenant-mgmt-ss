package tenantdbserv

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"os"

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
		log.Fatal("Failed to append PEM.")
	}

	mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})

	cfgdb := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBHOST"),
		DBName:               "tenants",
		AllowNativePasswords: true,
		TLSConfig:            "custom",
	}

	// Get a database handle.
	var err error
	tenantsDbHandle, err = sql.Open("mysql", cfgdb.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := tenantsDbHandle.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("connected to tenants db on server")
	return tenantsDbHandle
}

func (tdb *tenantDb) PingDb() error {
	err := tdb.Handle.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
