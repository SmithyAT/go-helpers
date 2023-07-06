package dbhelper

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"net"
	"time"
)

// PgLogDSN returns the DSN string for the MySQL database with the password masked
func PgLogDSN(dbConfig DbConfig) (dsn string) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbConfig.Username, "*****", dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// PgSQLConnect connects to the Postgres database and returns a pointer to the connection
func PgSQLConnect(dbConfig DbConfig) (db *sqlx.DB, disconnect func(), dbErr error) {
	err := connCheck(dbConfig.Host, dbConfig.Port)
	if err != nil {
		return nil, nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

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
func connCheck(host, port string) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second*5)
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
