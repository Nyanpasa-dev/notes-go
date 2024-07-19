package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

type User struct {
	gorm.Model
	Username string
	Password string
	Avatar   string
	Login    string
	Email    string
}
