package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

// the init function runs before main automatically
func init() {
	
	// if fly_app_name exists then we are in production, otherwise we load the local .env file
	if os.Getenv("FLY_APP_NAME") == "" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Warning: Failed to load .env file")
		}
	}

	var err error
	DB, err = SetupDB()
	if err != nil{
		// handle properly later
		log.Fatal("Connection to database failed")
	}

	RunSchema(DB)

}

func SetupDB() (*sql.DB, error) {
    dbKey := os.Getenv("DB_KEY")
    dbUrl := fmt.Sprintf("libsql://cms-3b00d09.turso.io?authToken=%s", dbKey)
    return sql.Open("libsql", dbUrl)
}