package models

type User struct {
	Token      int64
	Favourites []Note
}

type Note struct {
	City string
	Note string
}
