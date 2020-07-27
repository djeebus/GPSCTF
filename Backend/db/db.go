package db

import (
	"database/sql"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var sqlDb *sql.DB

func ensureDatabaseExists() string {
	dbPath := viper.GetString("DatabasePath")
	if len(dbPath) == 0 {
		dbPath = "gpsctf.db"
	}

	fullPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Failed to find abs path for %s: %v", dbPath, err)
	}

	// todo: this seem unnecessarily complex. really??
	info, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to find info on the db file: %v", err)
	} else if info != nil && info.IsDir() {
		log.Fatalf("Database path '%s' is a directory", fullPath)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("Failed to create '%s': %v", fullPath, err)
	}

	defer file.Close()

	return fullPath
}

func OpenDatabase() {
	var err error

	dbPath := ensureDatabaseExists()
	log.Printf("Database path: %s", dbPath)

	sqlDb, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open '%s': %v", dbPath, err)
	}

	tables, err := getTableList()
	if err != nil {
		log.Fatalf("Failed to get table list: %v", err)
	}

	if _, ok := tables["games"]; !ok {
		err = createTable(gamesTable)
		if err != nil {
			log.Fatalf("Failed to create games table: %v", err)
		}
	}

	if _, ok := tables["players"]; !ok {
		err = createTable(playersTable)
		if err != nil {
			log.Fatalf("Failed to create players table: %v", err)
		}
	}
}

func createTable(query string) error {
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func getTableList() (map[string]bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table'`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer row.Close()

	tables := make(map[string]bool)

	for row.Next() {
		var name string
		err = row.Scan(&name)
		if err != nil {
			return nil, err
		}

		tables[name] = true
	}

	return tables, nil
}

func CloseDatabase() error {
	err := sqlDb.Close()
	if err != nil {
		return err
	}

	return nil
}
