package domain

import (
	"time"
)

type Address struct {
	City   string `json:"city"`
	State  string `json:"state"`
	Street string `json:"street"`
	Number int    `json:"number"`
}

type User struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Age                  int    `json:"age"`
	NRC                  int    `json:"nrc"`
	DateOfRegistration   time.Time `json:"date_of_registration"`
	Address              Address   `json:"address"`
}
