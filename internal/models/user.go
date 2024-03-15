package models

import "time"

type User struct {
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type FullUser struct {
	User       User
	Favourites []string
	Alerts     []Alert
}

type Alert struct {
	City string
	Date time.Time
}
