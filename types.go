package main

import "database/sql"

type (
	element struct {
		Title    string `json:"title"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	base struct {
		Rows []element
	}

	dataBase struct {
		postgres *sql.DB
	}

	response struct {
		Error     string    `json:"error"`
		Passwords []element `json:"passwords"`
	}
)
