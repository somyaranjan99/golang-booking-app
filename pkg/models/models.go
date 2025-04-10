package models

import (
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"time"
)

type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      reservationtypes.Room
}
