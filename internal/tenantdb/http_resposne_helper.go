package tenantdbserv

import (
	"database/sql"
	"fmt"
	"strings"

	subscripifylogger "dev.azure.com/Subscripify/subscripify-prod/subscripify-logger.git"
	"github.com/go-sql-driver/mysql"
)

func HttpResponseHelperSQLInsert(inner sql.Result, sqlErr error) (sql.Result, int, error) {

	if sqlErr != nil {
		//parse the sql error - add special cases here as needed
		me, ok := sqlErr.(*mysql.MySQLError)
		if !ok {
			return inner, 500, fmt.Errorf("db server error")
		}
		if me.Number == 1062 {
			subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
			return inner, 409, fmt.Errorf("fail on db insert")
		}
		if me.Number == 1452 {
			if strings.Contains(me.Message, "CONSTRAINT `fk_lord_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("lordServiceConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_public_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("publicServicesConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_super_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("superServiceConfig does not exist")
			}

		}
		subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
		return inner, 400, fmt.Errorf("fail on db insert")
	}

	return inner, 200, nil
}

func HttpResponseHelperSQLDelete(inner sql.Result, sqlErr error) (sql.Result, int, error) {

	if sqlErr != nil {
		//parse the sql error - add special cases here as needed
		me, ok := sqlErr.(*mysql.MySQLError)
		if !ok {
			return inner, 500, fmt.Errorf("db server error")
		}

		if me.Number == 1451 {

			subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
			return inner, 400, fmt.Errorf("can not delete, this tenant still has one or more vassal tenants")
		}
		subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
		return inner, 400, fmt.Errorf("fail on db delete")
	}
	return inner, 200, nil
}

func HttpResponseHelperSQLUpdate(inner sql.Result, sqlErr error) (sql.Result, int, error) {

	if sqlErr != nil {
		//parse the sql error - add special cases here as needed
		me, ok := sqlErr.(*mysql.MySQLError)
		if !ok {
			return inner, 500, fmt.Errorf("db server error")
		}
		if me.Number == 1452 {
			if strings.Contains(me.Message, "CONSTRAINT `fk_lord_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("lordServiceConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_public_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("publicServicesConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_super_services_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("superServiceConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_private_access_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("privateAccessConfig does not exist")
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_custom_access_config`") {
				subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
				return inner, 400, fmt.Errorf("customAccessConfig does not exist")
			}
		}
		subscripifylogger.DebugLog.Printf("sql command error: %s", sqlErr.Error())
		return inner, 400, fmt.Errorf("fail on db update")
	}
	return inner, 200, nil
}
