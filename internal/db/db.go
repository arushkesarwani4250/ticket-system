package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	if err := setupSchema(db); err != nil {
		return nil, fmt.Errorf("failed to set up schema: %w", err)
	}

	return db, nil
}

func setupSchema(db *sql.DB) error {
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	ticketsTableQuery := `
	CREATE TABLE IF NOT EXISTS tickets (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		"desc" TEXT,
		status VARCHAR(50) CHECK(status IN ('open', 'in_progress', 'closed')) DEFAULT 'open',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	userTicketsTableQuery := `
	CREATE TABLE IF NOT EXISTS user_tickets (
		user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
		ticket_id INTEGER REFERENCES tickets(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, ticket_id)
	);`

	if _, err := db.Exec(usersTableQuery); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	if _, err := db.Exec(ticketsTableQuery); err != nil {
		return fmt.Errorf("failed to create tickets table: %w", err)
	}
	if _, err := db.Exec(userTicketsTableQuery); err != nil {
		return fmt.Errorf("failed to create user_tickets table: %w", err)
	}

	log.Println("Database schema updated")
	return nil
}
