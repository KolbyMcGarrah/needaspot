package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var Db *sql.DB

// InitDB creates a connection to the database
func InitDB() {
	// set up db configs
	const (
		host     = "localhost"
		port     = 5432
		user     = "nas_user"
		password = "na5ty1"
		dbname   = "nasdb"
	)
	// connect to the database
	pqInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", pqInfo)
	//log errors if any occur
	if err != nil {
		log.Panic(err)
	}

	//test the connection
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db

}

// Migrate runs migration files for us.
func Migrate() {
	//check db connection
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	driver, _ := postgres.WithInstance(Db, &postgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/postgres",
		"postgres",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}
