package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/subscripify-logger.git"
	tenant "dev.azure.com/Subscripify/subscripify-prod/tenant-mgmt-ss.git/cmd/tenant"
	tenantapi "dev.azure.com/Subscripify/subscripify-prod/tenant-mgmt-ss.git/internal/httpclient"
)

func Execute() {
	testDataCmd := flag.NewFlagSet("test-data", flag.ExitOnError)
	testDataGen := testDataCmd.Bool("gen", false, "generate")
	testDataDel := testDataCmd.Bool("del", false, "delete")
	testDataPrg := testDataCmd.Bool("prg", false, "purge")

	if len(os.Args) < 2 {
		fmt.Println("subcommand required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "test-data":
		testDataCmd.Parse(os.Args[2:])
		env := os.Getenv("SUBSCRIPIFY_DB_ENV")
		if env != "localdb" {
			subscripifylogger.FatalLog.Fatalf(`SUBSCRIPIFY_DB_ENV needs to be set to "localdb" to generate test data`)
		}
		if *testDataGen {

			tenant.TestDataCreate()
		}
		if *testDataDel {

			tenant.TestDataDelete()
		}
		if *testDataPrg {

			tenant.TestDataPurge()
		}
	case "serve":

		subscripifylogger.InfoLog.Printf("tenant management service started")
		router := tenantapi.NewRouter()
		log.Fatal(http.ListenAndServe(":8080", router))
	default:
		subscripifylogger.InfoLog.Printf(`invalid subcommand used :"%s" disconnecting database and closing application`, os.Args[1])
	}

}
