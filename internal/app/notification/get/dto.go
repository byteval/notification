package get

import (
	"time"
)

type Request struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type Response struct {
	ID        string    `json:"id"`
	Title     string    `json:"subject"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
