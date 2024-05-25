package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func NewCustomerDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./customer.db")
	if err != nil {
		slog.Error("Error while attempting to open database", "error", err)
		panic("Error while attempting to open database")
	}

	err = db.Ping()

	if err != nil {
		slog.Error("Error while attempting to ping database. Need to shutdown the application", "error", err)
		panic("Error while attempting to ping database")
	}

	slog.Info("Database connected", "db", db)
	return db
}

func InitDb() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./customer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the customers table
	createTableSQL := `CREATE TABLE IF NOT EXISTS customers (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" TEXT,
        "age" INTEGER
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v\n", err)
	}
	fmt.Println("Table created successfully!")
}

func SeedDb() {
	db, err := sql.Open("sqlite3", "./customer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new customer
	customer := Customer{
		ID:   0,
		Name: "do",
		Age:  28,
	}

	// Insert the customer into the table
	insertCustomerSQL := `INSERT INTO customers (name, age) VALUES (?, ?)`
	_, err = db.Exec(insertCustomerSQL, customer.Name, customer.Age)
	if err != nil {
		log.Fatalf("Error inserting customer: %v\n", err)
	}
	fmt.Println("Customer inserted successfully!")

}

// getCustomer retrieves a customer by name from the database
func getCustomer(db *sql.DB, name string) (Customer, error) {
	var customer Customer
	querySQL := `SELECT id, name, age FROM customers WHERE name = ?`
	row := db.QueryRow(querySQL, name)
	err := row.Scan(&customer.ID, &customer.Name, &customer.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return customer, fmt.Errorf("no customer found with name %s", name)
		}
		return customer, err
	}
	return customer, nil
}

type dbKey int

var dbCtxKey dbKey

func WithDbContext(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey, db)
}

func FromDbContext(ctx context.Context) (*sql.DB, bool) {
	db, ok := ctx.Value(dbCtxKey).(*sql.DB)
	return db, ok
}
