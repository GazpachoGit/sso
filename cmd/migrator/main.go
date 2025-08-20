package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	//driver for migrations for sqlite
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	//driver for migrations from files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "")
	flag.StringVar(&migrationsPath, "storage-path", "", "")
	flag.StringVar(&migrationsTable, "storage-path", "", "")

	if storagePath == "" {
		panic("empty storagePath")
	}
	if migrationsPath == "" {
		panic("empty migrationsPath")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
		}
		panic(err)
	}
	fmt.Println("migration applied succesfully")
}
