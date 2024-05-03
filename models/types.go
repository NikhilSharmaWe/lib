package models

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Book struct {
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
	Year   int    `json:"year,omitempty"`
}
