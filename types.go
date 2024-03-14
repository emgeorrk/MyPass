package main

type (
	Element struct {
		Title    string `json:"title"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	Base struct {
		Rows []Element
	}
)
