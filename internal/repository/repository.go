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
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedRservation(id int, processed int) error
	AllRooms() ([]models.Room, error)
	GetRoomRestrictionsByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	DeleteBlockById(id int) error
	InsertBlockForRoom(id int, startDate time.Time) error
}
