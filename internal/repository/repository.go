package repository

import (
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
)

type DatabaseRepo interface {
	AllUsers() bool
	BookReservation(*reservationtypes.Reservation) (string, error)
}
