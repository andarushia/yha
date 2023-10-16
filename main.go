package main

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Creds struct {
	Name       string
	Surname    string
	Patronymic string
	Age        uint64
	Sex        string
	Origin     string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "anya"
	password = "sqlxpass"
	dbname   = "anyatop"
)

func main() {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgresql", psql)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("connection established")

	rows, err := getNames(db)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}

func getNames(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM profile")
	return rows, err
}
