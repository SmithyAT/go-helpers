package dbhelper

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// MySqlLogDSN creates a masked version of the Data Source Name (DSN) for a MySQL database connection.
// The DSN string is created using the DbConfig object, only it masks the password for security reasons.
// The returned DSN string includes the username, masked password, host, port, and database name.
// Two additional parameters are added: a connection timeout of 5 seconds and automatic commit mode enabled.
// This function is especially useful when logging or displaying the DSN without exposing sensitive information.
//
// The function accepts a single parameter:
// - dbConfig: a DbConfig object that includes the database username, password, host, port, and database name.
//
// Returns:
// - dsn: a string representing the masked DSN of the MySQL database connection.
func MySqlLogDSN(dbConfig DbConfig) (dsn string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&autocommit=true",
		dbConfig.Username, "******", dbConfig.Host, dbConfig.Port, dbConfig.Database)
}

// MySQLConnect is a function that allows the application to establish a connection with the MySQL database.
// This function takes a DbConfig structure as an input parameter which includes the configurations required to connect like username, password, host, port and name of the database.
//
// It utilizes the sqlx library's Connect function to create a connection with the MySQL database. If the connection is established successfully,
// it sets the maximum connection lifetime, maximum number of open connections and maximum idle connections for the MySQL database.
//
// If a connection is not established, the function returns an error. Otherwise, it returns a pointer to the DB and a function to close connection when done.
//
// Usage:
//
//	db, disconnect, err := MySQLConnect(cfg)
//	if err != nil {
//	    log.Fatalln(err)
//	}
//	// Make sure to close the DB connection when you're done using it
//	defer disconnect()
//
// Parameters:
//
// dbConfig: Specifies the configuration required for connecting to MySQL.
//
// Returns:
//
// *sqlx.DB: Returns a SQLx DB object which can be used to interact with the database.
//
// disconnect: It's a function you can call to close the DB connection.
//
// error: If there is any error encountered during the process, it will be returned as an error type
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
