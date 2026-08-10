package importer

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func Connect() (*Database, error) {
	dsn := "opencart:opencart123@tcp(localhost:3306)/opencart?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to MariaDB")

	return &Database{
		DB: db,
	}, nil
}