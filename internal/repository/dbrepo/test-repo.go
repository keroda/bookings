package dbrepo

import (
	"errors"
	"time"

	"github.com/keroda/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) AddReservation(res models.Reservation) (int, error) {
	if res.RoomID == 3 {
		return 0, errors.New("add-reservation error")
	}
	return 1, nil
}

func (m *testDBRepo) AddRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID > 2 {
		return errors.New("add-restriction error")
	}
	return nil
}

// return true means room is available, else room is taken
func (m *testDBRepo) SearchAvailabilityByRoom(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Room error")
	}
	return room, nil
}
