package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host        = "localhost"
	port        = 5432
	user        = "egormerk"
	sqlPassword = "admin"
	dbname      = "postgres"
)

var db = new(dataBase)

func (e *element) getElem(key string, from *gin.Context) (int, error) {
	jsonElem := from.GetHeader(key)
	if jsonElem == "" {
		errorString := fmt.Sprintf("Error: header '%s' is empty", key)
		return http.StatusBadRequest, errors.New(errorString)
	}

	err := json.Unmarshal([]byte(jsonElem), &e)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}

func (d *dataBase) getBase() (*base, int, error) {
	row, err := d.postgres.Query("SELECT COUNT(*) count FROM passwords")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer row.Close()

	var count int
	for row.Next() {
		if err := row.Scan(&count); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	if count == 0 {
		return nil, http.StatusBadRequest, errors.New("table is empty")
	}

	elements := make([]element, 0)
	rows, err := d.postgres.Query("SELECT title, login, password FROM passwords")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, login, pass string
		if err := rows.Scan(&title, &login, &pass); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		elements = append(elements, element{title, login, pass})
	}

	return &base{elements}, 0, nil
}

func (d *dataBase) addElem(newElem element) (int, error) {
	row, err := d.postgres.Query("SELECT COUNT(title) count FROM passwords WHERE title = $1", newElem.Title)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer row.Close()

	var count int
	row.Next()
	if err := row.Scan(&count); err != nil {
		return http.StatusInternalServerError, err
	}
	if count != 0 {
		return http.StatusBadRequest, errors.New("element already exists")
	}

	_, err = d.postgres.Exec("INSERT INTO passwords (title, login, password) VALUES ($1, $2, $3)",
		newElem.Title, newElem.Login, newElem.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

func (d *dataBase) editElem(oldElem, newElem element) (int, error) {
	row, err := d.postgres.Query(
		"SELECT COUNT(title) count FROM passwords WHERE title = $1 AND login = $2 AND password = $3",
		oldElem.Title, oldElem.Login, oldElem.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer row.Close()

	var count int
	row.Next()
	if err := row.Scan(&count); err != nil {
		return http.StatusInternalServerError, err
	}
	switch {
	case count < 1:
		return http.StatusBadRequest, errors.New("element doesn't exists")
	case count > 1:
		return http.StatusInternalServerError, errors.New("element present more than once")
	}

	_, err = d.postgres.Exec(
		"UPDATE passwords SET title=$1, login=$2, password=$3 WHERE title=$4 AND login=$5 AND password=$6",
		newElem.Title, newElem.Login, newElem.Password,
		oldElem.Title, oldElem.Login, oldElem.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (d *dataBase) removeElem(elem element) (int, error) {
	row, err := d.postgres.Query(
		"SELECT COUNT(title) count FROM passwords WHERE title = $1 AND login = $2 AND password = $3",
		elem.Title, elem.Login, elem.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer row.Close()

	var count int
	row.Next()
	if err := row.Scan(&count); err != nil {
		return http.StatusInternalServerError, err
	}
	switch {
	case count < 1:
		return http.StatusBadRequest, errors.New("element doesn't exists")
	case count > 1:
		return http.StatusInternalServerError, errors.New("element present more than once")
	}

	_, err = d.postgres.Exec(
		"DELETE FROM passwords WHERE title=$1 AND login=$2 AND password=$3",
		elem.Title, elem.Login, elem.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (d *dataBase) connectDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, sqlPassword, dbname)
	fmt.Println(psqlInfo)

	var err error
	d.postgres, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = d.postgres.Ping()
	if err != nil {
		return err
	}

	_, err = d.postgres.Exec("CREATE TABLE IF NOT EXISTS passwords (id SERIAL PRIMARY KEY, title VARCHAR, login VARCHAR, password VARCHAR);")
	if err != nil {
		return err
	}

	log.Println("Successfully connected to the database")
	return nil
}
