package entity

import "time"

type Dollar struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type DollarResponse struct {
	Currency struct {
		Value string `json:"bid"`
	} `json:"USDBRL"`
}
