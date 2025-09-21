// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	// The underscore before the package name means we are importing it
	// for its side effects, which in this case is to register the
	// odbc driver with the database/sql package.
	_ "github.com/alexbrainman/odbc"
)

func main() {
	// Define the DSN (Data Source Name) you want to connect to.
	// In your case, it's "CLEARSCAPE".
	dsn := "CLEARSCAPE"

	// The connection string for the ODBC driver.
	// The format is "dsn=<YOUR_DSN_NAME>".
	connectionString := fmt.Sprintf("dsn=%s", dsn)

	fmt.Printf("Attempting to connect to DSN: %s\n", dsn)

	// Open a connection to the database.
	// The first argument is the driver name, which is "odbc" for this package.
	// The second argument is the connection string.
	db, err := sql.Open("odbc", connectionString)
	if err != nil {
		// If sql.Open fails, it will return an error.
		// We use log.Fatalf to print the error and exit the program.
		log.Fatalf("Error creating connection to database: %v", err)
	}

	// It's crucial to close the connection when you're done with it
	// to free up resources. The defer statement ensures this runs
	// right before the main function exits.
	defer db.Close()

	// To verify the connection is actually alive, we can ping the database.
	// Ping() will establish a connection if one hasn't been established yet.
	err = db.Ping()
	if err != nil {
		// If Ping() returns an error, it means we couldn't connect.
		log.Fatalf("Error connecting to database: %v", err)
	}

	// If we reach this point, the connection was successful!
	fmt.Println("Successfully connected to the CLEARSCAPE DSN! 🎉")

	// You can now proceed to execute queries, etc.
	// For example:
	rows, err := db.Query("SELECT top 10 * FROM DBC.tablesv")
	// ...
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	// ... Process the rows here ...
	fmt.Println("Query executed successfully!")
}
