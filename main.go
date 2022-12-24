package main

import (
	"database/sql"

	"github.com/Trickster-ID/dbo/config"
)

var (
	db *sql.DB = config.SetUpDatabaseConnection()
)

func main() {
	defer config.CloseDatabaseConnection(db)

}
