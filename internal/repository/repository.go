package repository

import (
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	BookReservation(*reservationtypes.Reservation) (string, error)
	SearchAvailabilityForAllRooms(start_date, end_date time.Time) ([]reservationtypes.Room, error)
	GetRoomByID(id int) (reservationtypes.Room, error)
}
