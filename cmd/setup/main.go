package main

import (
	"database/sql"
	"fmt"
	"log"
	"shop-api-go/internal/core/domain"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// input contains all information need to add a new admin.
type input struct {
	adminUsername string
	adminEmail    string
	adminPassword string
	databaseURL   string
}

// readInput reads the input from os.Stdin.
func readInput() (*input, error) {
	fmt.Println("Enter admin username: ")
	var username string
	_, err := fmt.Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("reading admin username: %v", err)
	}
	username = strings.TrimSpace(username)

	fmt.Println("Enter admin email: ")
	var email string
	_, err = fmt.Scan(&email)
	if err != nil {
		return nil, fmt.Errorf("reading admin email: %v", err)
	}
	email = strings.TrimSpace(email)

	fmt.Println("Enter admin password: ")
	var password string
	_, err = fmt.Scan(&password)
	if err != nil {
		return nil, fmt.Errorf("reading admin password: %v", err)
	}
	password = strings.TrimSpace(password)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %v", err)
	}

	fmt.Println("Enter database URL: ")
	var databaseURL string
	_, err = fmt.Scan(&databaseURL)
	if err != nil {
		return nil, fmt.Errorf("reading database URL: %v", err)
	}
	databaseURL = strings.TrimSpace(databaseURL)

	return &input{
		adminUsername: username,
		adminEmail:    email,
		adminPassword: string(hash),
		databaseURL:   databaseURL,
	}, nil
}

func main() {
	rInput, err := readInput()
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	db, err := sql.Open("postgres", rInput.databaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	if db.Ping() != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO users(id, username, email, password, role)
		VALUES($1, $2, $3, $4, $5)`,
		uuid.New(), rInput.adminUsername, rInput.adminEmail, rInput.adminPassword, domain.Admin,
	)

	if err != nil {
		log.Fatalf("error inserting user: %v", err)
	}

	fmt.Println("Successfully created first admin user")
}
