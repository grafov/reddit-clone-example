package storage

import (
	"database/sql"
	"errors"

	"redditclone/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grafov/kiwi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	// Database initialized in this package and then we can use DB
	// methods everywhere. Note that "repository" pattern is
	// redundant here because we don't change the initialized once
	// instance.
	DB *sqlx.DB
	log = kiwi.Fork().With("service", "database")
)

func Init() {
	db, err := sql.Open("postgres", config.DB.ConnectionString())
	if err != nil {
		log.Log("info", "opening connection", "err", err)
		return
	}
	if _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + config.DB.Namespace); err != nil {
		log.Log("info", "creating schema", "err", err)
		return
	}

	log.Log("info", "SET search_path="+config.DB.Namespace)
	if _, err = db.Exec("SET search_path=" + config.DB.Namespace); err != nil {
		log.Log("info", "creating schema", "err", err)
		return
	}
	if err = rolling(db); err != nil {
		log.Log("info", "apply migration", "config", config.DB, "err", err)
		return
	}
	if config.DB.MigrateDown > 0 {
		log.Log("info", "migration rolled down")
		return
	}
	if _, err = db.Exec("SET search_path=" + config.DB.Namespace); err != nil {
		log.Log("info", "creating schema", "err", err)
		return
	}
	DB = sqlx.NewDb(db, "postgres")
}

// migrate create -dir database/migrations -ext sql NAME
func rolling(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.New("database instance assertion:" + err.Error())
	}
	instance, err := migrate.NewWithDatabaseInstance(
		"file://storage/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return errors.New("migrations instance assertion:" + err.Error())
	}
	if config.DB.MigrateDown > 0 {
		if err = instance.Steps(-1 * config.DB.MigrateDown); err != nil {
			return errors.New("migrations roll back:" + err.Error())
		}
		return nil
	}
	if err = instance.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.New("migration rolling up:" + err.Error())
	}

	return nil
}
