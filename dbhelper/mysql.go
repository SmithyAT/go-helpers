package dbhelper

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// MySqlLogDSN returns the DSN string for the MySQL database with the password masked
func MySqlLogDSN(dbConfig DbConfig) (dsn string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&autocommit=true",
		dbConfig.Username, "******", dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// MySQLConnect connects to the MySQL database and returns a pointer to the connection
func MySQLConnect(dbConfig DbConfig) (db *sqlx.DB, disconnect func(), dbErr error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&autocommit=true",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

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
