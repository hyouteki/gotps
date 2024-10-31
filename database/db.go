package database

import (
	"database/sql"
	"log"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Connect(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("error: failed to open database `%s`: %v\n", dbPath, err)
		return err
	}
	return nil
}

func Init(sqlFilePath string) error {
	content, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Printf("error: failed to read file `%s`: %v\n", sqlFilePath, err)
		return err
	}

	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("error: failed to execute query `%s`: %v\n", query, err)
			return err
		}
	}
	return nil
}
