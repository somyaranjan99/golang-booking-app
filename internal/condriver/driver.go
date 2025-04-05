package condriver

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	SQL *sql.DB
}

func Dbconn() (*DB, error) {
	// Connection string format: "username:password@tcp(host:port)/dbname"
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/bookings?parseTime=true")
	if err != nil {
		log.Println("Failed to open MySQL connection:", err)
		return nil, err
	}

	// Verify the connection with Ping()
	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping MySQL:", err)
		db.Close() // Close if Ping fails
		return nil, err
	}

	log.Println("Successfully connected to MySQL!")
	return &DB{SQL: db}, nil
}
