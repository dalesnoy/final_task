package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbfile string) error {

	var install bool
	if _, err := os.Stat(dbfile); err != nil {
		install = true
	}

	var err error
	if DB, err = sql.Open("sqlite", dbfile); err != nil {
		return err
	}

	if install {
		if _, err := DB.Exec(`
			CREATE TABLE scheduler (
                id      INTEGER PRIMARY KEY AUTOINCREMENT,
                date    CHAR(8)      NOT NULL DEFAULT "",
                title   VARCHAR(256) NOT NULL DEFAULT "",
                comment TEXT         NOT NULL DEFAULT "",
                repeat  VARCHAR(128) NOT NULL DEFAULT ""
            );
            CREATE INDEX idx_date ON scheduler(date);
        `); err != nil {
			return err
		}
	}

	return nil
}
