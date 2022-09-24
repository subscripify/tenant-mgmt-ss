package tenantdbserv

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func InsertResponseHelper(inner sql.Result, sqlErr error) (r sql.Result, httpResponseCode int, message string) {

	if sqlErr != nil {
		//parse the sql error - add special cases here as needed
		me, ok := sqlErr.(*mysql.MySQLError)
		if !ok {
			return inner, 500, "server error"
		}
		if me.Number == 1062 {
			return inner, 409, fmt.Sprintf("fail on db insert: %s", sqlErr.Error())
		}
		if me.Number == 1452 {
			if strings.Contains(me.Message, "CONSTRAINT `fk_lord_services_config`") {
				return inner, 400, fmt.Sprintf("lordServiceConfig does not exist: %s", sqlErr.Error())
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_public_services_config`") {
				return inner, 400, fmt.Sprintf("publicServicesConfig does not exist: %s", sqlErr.Error())
			}
			if strings.Contains(me.Message, "CONSTRAINT `fk_super_services_config`") {
				return inner, 400, fmt.Sprintf("superServiceConfig does not exist: %s", sqlErr.Error())
			}

		}
		return inner, 400, fmt.Sprintf("fail on db insert: %s", sqlErr.Error())
	}

	return inner, 200, "successful insert"
}
