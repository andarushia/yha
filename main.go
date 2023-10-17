package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type Creds struct {
	Id         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        *uint64
	Sex        *string
	Origin     *string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "anya"
	password = "sqlxpass"
	dbname   = "anyadb"
)

func main() {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password,
		host, port, dbname)
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		handleErr("unable to connect to database", err)
	}
	defer conn.Close(context.Background())

	err = getName(conn)
	if err != nil {
		handleErr("couldn't get rows", err)
	}
}

func getName(conn *pgx.Conn) error {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM public.profile")
	people, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Creds])
	if err != nil {
		return err
	}

	return nil
}

func handleErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", msg, err)
		os.Exit(1)
	}
}
