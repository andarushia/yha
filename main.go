package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

type Age struct {
	Age uint8 `json:"age"`
}

func (str Age) getJson(requestUrl string) uint8 {
	data, _ := http.Get(requestUrl)
	responseBody, _ := io.ReadAll(data.Body)
	err := json.Unmarshal(responseBody, &str)
	if err != nil {
		handleErr("error during requesting data", err)
	}
	return str.Age
}

type Gender struct {
	Sex string `json:"gender"`
}

func (str Gender) getJson(requestUrl string) string {
	data, _ := http.Get(requestUrl)
	responseBody, _ := io.ReadAll(data.Body)
	err := json.Unmarshal(responseBody, &str)
	if err != nil {
		handleErr("error during requesting data", err)
	}
	return str.Sex
}

type Country struct {
	Origin []struct {
		Country string `json:"country_id"`
	} `json:"country"`
}

func (str Country) getJson(requestUrl string) string {
	data, _ := http.Get(requestUrl)
	responseBody, _ := io.ReadAll(data.Body)
	err := json.Unmarshal(responseBody, &str)
	if err != nil {
		handleErr("error during requesting data", err)
	}
	return str.Origin[0].Country
}

func requestData(wg sync.WaitGroup, ch chan<- []byte, requestUrl string) {

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
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		handleErr("unable to connect to database", err)
	}
	defer dbpool.Close()

	cliPrompt(dbpool)

	err = getNameAndPopulate(dbpool)
	if err != nil {
		handleErr("couldn't get rows", err)
	}
}

func cliPrompt(dbpool *pgxpool.Pool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("select:\n1. add entry\n2. delete entry\n3. print table\n")
	fmt.Print("your input: ")
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		handleErr("error parsing data", err)
	}
	switch choice {
	case 1:
		fmt.Println("provide name, surname and patronymic if exists:")
		promptString, err := reader.ReadString('\n')
		if err != nil {
			handleErr("error parsing data", err)
		}
		name := strings.Split(promptString, " ")
		if len(name) == 2 {
			addEntry(dbpool, name[0], name[1], "")
		} else if len(name) == 3 {
			addEntry(dbpool, name[0], name[1], name[2])
		} else {
			handleErr("yhahha", nil)
		}
	case 2:
		fmt.Println("not implemented yet")
	case 3:
		fmt.Println("not implemented yet")
	}
}

func getNameAndPopulate(dbpool *pgxpool.Pool) error {
	rows, _ := dbpool.Query(context.Background(), "SELECT * FROM public.profile")
	var person Creds
	_, err := pgx.ForEachRow(rows, []any{&person.Id, &person.Name, &person.Surname,
		&person.Patronymic, &person.Age, &person.Sex, &person.Origin}, func() error {

		age, sex, origin := populate(person.Name)
		replaceQuery(dbpool, person.Id, age, sex, origin)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func handleErr(errorMsg string, err error) {
	fmt.Fprintf(os.Stderr, "%v: %v\n", errorMsg, err)
	os.Exit(1)
}

func populate(name string) (uint8, string, string) {
	var (
		age    Age
		sex    Gender
		origin Country
	)

	agifyUrl := fmt.Sprintf("https://api.agify.io/?name=%v", name)
	genderizeUrl := fmt.Sprintf("https://api.genderize.io/?name=%v", name)
	nationalizeUrl := fmt.Sprintf("https://api.nationalize.io/?name=%v", name)

	return age.getJson(agifyUrl), sex.getJson(genderizeUrl), origin.getJson(nationalizeUrl)
}

func replaceQuery(dbpool *pgxpool.Pool, id uint64, age uint8, sex string, origin string) {
	queryString := fmt.Sprintf("UPDATE profile SET age = %v, sex = '%v', origin = '%v' WHERE id = %v;", age, sex, origin, id)
	_, err := dbpool.Exec(context.Background(), queryString)
	if err != nil {
		handleErr("error while updating database", err)
	}
}

func addEntry(dbpool *pgxpool.Pool, name string, surname string, patronymic string) {
	var queryString string
	if patronymic != "" {
		queryString = fmt.Sprintf("INSERT INTO profile (name, surname, patronymic) VALUES ('%v', '%v', '%v');", name, surname, patronymic)
	} else {
		queryString = fmt.Sprintf("INSERT INTO profile (name, surname) VALUES ('%v', '%v');", name, surname)
	}
	_, err := dbpool.Exec(context.Background(), queryString)
	if err != nil {
		handleErr("error while inserting data", err)
	}
	searchString := fmt.Sprintf("SELECT id FROM public.profile WHERE name = %v AND surname = %v", name, surname)
	var person Creds
	dbpool.QueryRow(context.Background(), searchString).Scan(&person)
	age, sex, origin := populate(name)
	replaceQuery(dbpool, person.Id, age, sex, origin)
}
