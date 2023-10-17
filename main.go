package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
	"os"
)

type Creds struct {
	Id         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        *uint8
	Sex        *string
	Origin     *string
}

type Message struct {
	Age    uint8  `json:"age"`
	Sex    string `json:"gender"`
	Origin []struct {
		Country string `json:"country_id"`
	} `json:"country"`
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

	err = getNameAndPopulate(conn)
	if err != nil {
		handleErr("couldn't get rows", err)
	}
}

func getNameAndPopulate(conn *pgx.Conn) error {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM public.profile")
	var person Creds
	_, err := pgx.ForEachRow(rows, []any{&person.Id, &person.Name, &person.Surname,
		&person.Patronymic, &person.Age, &person.Sex, &person.Origin}, func() error {

		msg := Populate(person.Name)

		fmt.Println(msg)

		return nil
	})
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

func Populate(name string) Message {
	var msg Message

	agifyUrl := fmt.Sprintf("https://api.agify.io/?name=%v", name)
	genderizeUrl := fmt.Sprintf("https://api.genderize.io/?name=%v", name)
	nationalizeUrl := fmt.Sprintf("https://api.nationalize.io/?name=%v", name)

	getJson(agifyUrl, msg)
	getJson(genderizeUrl, msg)
	getJson(nationalizeUrl, msg)

	return msg
}

func getJson(requestUrl string, msg Message) {
	data, _ := http.Get(requestUrl)
	responseBody, _ := io.ReadAll(data.Body)
	fmt.Println(responseBody)
	err := json.Unmarshal(responseBody, &msg)
	if err != nil {
		handleErr("error during requesting data", err)
	}
	fmt.Println(msg)
}
