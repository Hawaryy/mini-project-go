package migration

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	// Create users table
	if err := createUsersTable(db); err != nil {
		return err
	}

	// Create categories table
	if err := createCategoriesTable(db); err != nil {
		return err
	}

	// Create books table
	if err := createBooksTable(db); err != nil {
		return err
	}

	// Insert default user
	if err := insertDefaultUser(db); err != nil {
		return err
	}

	return nil
}

func createUsersTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100),
		modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		modified_by VARCHAR(100)
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}
	return nil
}

func createCategoriesTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100),
		modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		modified_by VARCHAR(100)
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating categories table: %v", err)
	}
	return nil
}

func createBooksTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		image_url VARCHAR(255),
		release_year INTEGER,
		price INTEGER,
		total_page INTEGER,
		thickness VARCHAR(50),
		category_id INTEGER NOT NULL REFERENCES categories(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100),
		modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		modified_by VARCHAR(100)
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating books table: %v", err)
	}
	return nil
}

func insertDefaultUser(db *sql.DB) error {
	// Check if user already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", "admin").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // User already exists
	}

	// Hash password (in production, use proper hashing)
	// For demo: admin123 -> will be hashed
	insertSQL := `
	INSERT INTO users (username, password, created_by, modified_by)
	VALUES ($1, $2, $3, $4)
	`
	// Password: admin123 (SHA256 hashed for demo)
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36icoegS" // bcrypt hash of "admin123"
	_, err = db.Exec(insertSQL, "admin", hashedPassword, "system", "system")
	if err != nil {
		return fmt.Errorf("error inserting default user: %v", err)
	}

	return nil
}
