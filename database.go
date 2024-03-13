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

type (
	Element struct {
		Title    string `json:"title"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	Base struct {
		Passwords []Element
	}
)

func getBase(db *sql.DB) (*Base, int, error) {
	row, err := db.Query("SELECT COUNT(*) count FROM passwords")
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
	rows, err := db.Query("SELECT title, login, password FROM passwords")
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

func addElem(db *sql.DB, title, login, password string) (int, error) {
	row, err := db.Query("SELECT COUNT(title) count FROM passwords WHERE title = $1", title)
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

	_, err = db.Exec("INSERT INTO passwords (title, login, password) VALUES ($1, $2, $3)",
		string(title), string(login), string(password))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func editElem(db *sql.DB, title, login, password string) error {

	return nil
}

func removeElem() {}

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, sqlPassword, dbname)
	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil

	//login := "EGor"
	//rows, err := db.Query("SELECT COUNT(title) count FROM passwords WHERE title = $1", string(login))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer rows.Close()
	//
	//// Проходимся по результатам запроса
	//for rows.Next() {
	//	var cnt int
	//	if err := rows.Scan(&cnt); err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("count: %d\n", cnt)
	//}
	//// Обработка возможной ошибки после вызова rows.Next()
	//if err := rows.Err(); err != nil {
	//	panic(err)
	//}
}
