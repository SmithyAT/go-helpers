package dbhelper

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"net"
	"regexp"
	"strings"
	"time"
)

// PgSqlConfig is the configuration for the MySQL database
type PgSqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// PgSqlDSN returns the DSN string for the MySQL database
func PgSqlDSN(dbConfig PgSqlConfig) (dsn string) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// PgSqlDSNMaskPassword returns the DSN string for the Postgres database with the password masked
func PgSqlDSNMaskPassword(s string) string {
	re, err := regexp.Compile(`(.*?):(.*?)@`)
	if err != nil {
		return ""
	}

	return re.ReplaceAllStringFunc(s, func(str string) string {
		parts := strings.Split(str, ":")
		return fmt.Sprintf("%s:%s@", parts[0], strings.Repeat("*", len(parts[1])-1)) // We subtract 1 from length to account for the "@" character
	})
}

// PgSQLConnect connects to the Postgres database and returns a pointer to the connection
func PgSQLConnect(dbConfig PgSqlConfig) (db *sqlx.DB, disconnect func(), dbErr error) {
	err := connCheck(dbConfig)
	if err != nil {
		return nil, nil, err
	}

	dsn := PgSqlDSN(dbConfig)

	pgDB, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, nil, err
	}

	return pgDB, func() {
		_ = pgDB.Close()
	}, nil
}

// connCheck make a quick tcp connection test ot the database host:port to see if it is available of not
// This is a workaround because the postgres driver has no timeout parameter and needs 150s to throw an error
func connCheck(dbConfig PgSqlConfig) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(dbConfig.Host, dbConfig.Port), time.Second*5)
	if err != nil {
		return err
	}
	if conn != nil {
		defer func(conn net.Conn) {
			_ = conn.Close()
		}(conn)
	}
	return nil
}
