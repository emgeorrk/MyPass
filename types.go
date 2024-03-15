package main

import "database/sql"

type (
	Element struct {
		Title    string `json:"title"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	base struct {
		Rows []Element
	}

	dataBase struct {
		postgres *sql.DB
	}
)
