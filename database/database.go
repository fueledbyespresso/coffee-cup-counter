package database

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB struct {
	Db           *sql.DB
	SessionStore *pgstore.PGStore
}

// InitDBConnection Initialize a database connection using the environment variable DATABASE_URL
//Returns type *sql.DB
func InitDBConnection() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Println("Cannot open SQL connection")
		panic(err.Error())
	}

	return db
}

// PerformMigrations Check that database is up to date.
//Will cycle through all changes in db/migrations until the database is up to date
func PerformMigrations(migrationsFolder string) {
	m, err := migrate.New(
		migrationsFolder,
		os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	fmt.Println("Database migrations completed. Database should be up to date")
}
