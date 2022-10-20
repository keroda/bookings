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
		return room, errors.New("room error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	if id > 2 {
		return u, errors.New("room error")
	}
	return u, nil
}
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

func (m *testDBRepo) UpdateReservation(r models.Reservation) error {

	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

func (m *testDBRepo) UpdateProcessedRservation(id int, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomRestrictionsByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

func (m *testDBRepo) DeleteBlockById(id int) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// sql := `
	// 	DELETE FROM room_restrictions WHERE id = $1
	// `
	// _, err := m.DB.Exec(ctx, sql, id)
	// if err != nil{
	// 	log.Println(err)
	// 	return err
	// }

	return nil
}
