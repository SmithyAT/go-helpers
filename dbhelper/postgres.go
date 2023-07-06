package dbhelper

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"net"
	"time"
)

// PgLogDSN constructs a Postgres Data Source Name (DSN) from the given DbConfig.
// However, for security reasons, it intentionally obscures the password
// within the DSN using a placeholder. It returns the DSN as a string.
//
// The DbConfig struct should contain information like Username, Host, Port,
// and Database name.
//
// The format of the returned DSN string is:
// "postgres://username:*****@host:port/database"
//
// Usage Example:
//
// cfg := DbConfig{Username: "user1", Host: "localhost", Port: "5432", Database: "mydb"}
// dsn := PgLogDSN(cfg)
// fmt.Println(dsn)  // Output: "postgres://user1:*****@localhost:5432/mydb"
func PgLogDSN(dbConfig DbConfig) (dsn string) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbConfig.Username, "*****", dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// PgSQLConnect takes a DbConfig structure as input and establishes a connection to
// a PostgreSQL database using the details provided in the provided configuration.
// If successful, the function returns a pointer to the sqlx.DB object representing
// the database connection, a function that can be called to disconnect from the
// database, and a nil error.
// If unsuccessful (e.g., because the connection could not be established or the
// connection check failed), the function returns a nil pointer, a nil disconnect
// function, and an error describing the failure.
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
