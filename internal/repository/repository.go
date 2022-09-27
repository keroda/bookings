package repository

import (
	"time"

	"github.com/keroda/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	AddReservation(res models.Reservation) (int, error)
	AddRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByRoom(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityAllRooms(start, end time.Time) ([]models.Room, error)
}
