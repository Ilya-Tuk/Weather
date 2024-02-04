package models

type User struct {
	Name       string
	Favourites []Note
}

type Note struct {
	City string
	Note string
}
