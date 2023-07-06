package dbhelper

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
	"time"
)

// MySqlConfig is the configuration for the MySQL database
type MySqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// MySqlDSN returns the DSN string for the MySQL database
func MySqlDSN(dbConfig MySqlConfig) (dsn string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&autocommit=true",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// MySqlDSNMaskPassword returns the DSN string for the MySQL database with the password masked
func MySqlDSNMaskPassword(s string) string {
	re, err := regexp.Compile(`(.*?):(.*?)@`)
	if err != nil {
		return ""
	}

	return re.ReplaceAllStringFunc(s, func(str string) string {
		parts := strings.Split(str, ":")
		return fmt.Sprintf("%s:%s@", parts[0], strings.Repeat("*", len(parts[1])-1)) // We subtract 1 from length to account for the "@" character
	})
}

// MySQLConnect connects to the MySQL database and returns a pointer to the connection
func MySQLConnect(dsn string) (db *sqlx.DB, disconnect func(), dbErr error) {
	mysqlDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}

	mysqlDB.SetConnMaxLifetime(time.Minute * 3)
	mysqlDB.SetMaxOpenConns(5)
	mysqlDB.SetMaxIdleConns(5)

	return mysqlDB, func() {
		_ = mysqlDB.Close()
	}, nil
}
