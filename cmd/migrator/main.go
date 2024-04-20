package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Storage path")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Migrator path")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Migrations table")

	flag.Parse()

	fmt.Println(storagePath, migrationsPath, migrationsTable)

	if storagePath == "" || migrationsPath == "" {
		fmt.Println("qwe", storagePath, migrationsPath)
		panic("storage-path and migrations-path are both required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")
}
