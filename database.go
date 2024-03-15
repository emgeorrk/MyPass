package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func (d *DataBase) getBase() (*Base, int, error) {
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

	elements := make([]Element, 0)
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
		elements = append(elements, Element{title, login, pass})
	}

	return &Base{elements}, http.StatusOK, nil
}

func (d *DataBase) addElem(newElem Element) (int, error) {
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
	return http.StatusOK, nil
}

func (d *DataBase) editElem(oldElem, newElem Element) (int, error) {
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

	return http.StatusOK, nil
}

func (d *DataBase) removeElem(elem Element) (int, error) {
	return http.StatusOK, nil
}

func (d *DataBase) connectDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, sqlPassword, dbname)
	fmt.Println(psqlInfo)

	var err error
	d.postgres, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	err = d.postgres.Ping()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Successfully connected to the database")
	return nil
}
