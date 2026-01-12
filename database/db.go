package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// SQL Command ran in pgAdmin
// --CREATE DATABASE mydb;
// --CREATE USER appuser WITH PASSWORD '0000';
// --GRANT ALL PRIVILEGES ON DATABASE mydb TO appuser;

//NOTE:
// Using database/sql along with pgx is an adapter pattern

const dbtype string = "postgres"

// const username string = "appuser"
const username string = "postgres"
const password string = "0000"
const dbhost string = "localhost"
const port string = "5432"
const dbname string = "mydb"
const security string = "sslmode=disable"

const DSN string = dbtype + "://" + username + ":" + password + "@" + dbhost + ":" + port + "/" + dbname + "?" + security

var DB *sql.DB
var err error

func Init() {
	DB, err = sql.Open("pgx", DSN)

	if err != nil {
		panic(fmt.Errorf("Cannot open up the databse %w", err))
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {

	createUserTable()
	createFileTable()

	fmt.Println("Db connected and table created")
}

func createUserTable() {
	query := `CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,          -- auto increment
				email VARCHAR(255) NOT NULL UNIQUE,
				first_name VARCHAR(255) NOT NULL,
				last_name VARCHAR(255) NOT NULL,
				password_hash TEXT NOT NULL,
				profile_pic TEXT,                  -- store URL or path
				created_at TIMESTAMPTZ DEFAULT now(),
				updated_at TIMESTAMPTZ DEFAULT now()
			);`

	_, err = DB.Exec(query)

	if err != nil {
		panic(fmt.Errorf("Cannot create the users Table %w", err))
	}
}

func createFileTable() {
	query := `CREATE TABLE IF NOT EXISTS files (
				id BIGSERIAL PRIMARY KEY,
				file_name TEXT NOT NULL,
				file_size BIGINT NOT NULL,
				mime_type TEXT NOT NULL,
				storage_key TEXT NOT NULL,                  
				created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, --ON DELETE CASCADE → if a user is deleted, their files are deleted automatically
				created_at TIMESTAMPTZ DEFAULT now(),
				updated_at TIMESTAMPTZ DEFAULT now()
			);`

	_, err = DB.Exec(query)

	if err != nil {
		panic(fmt.Errorf("Cannot create the files Table %w", err))
	}
}
