package tenant

import (
	"database/sql"
)

func updateRecord(value interface{}, field string, db *sql.DB) {
	// var t tenant
	// rows, err := db.Query("SELECT count(tenantUUID) where ? = ?", field, value)
}
