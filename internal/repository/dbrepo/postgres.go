package dbrepo

import (
	"context"
	"time"

	"github.com/keroda/bookings/internal/models"
)

func (m *postgresRepo) AllUsers() bool {
	return true
}

func (m *postgresRepo) AddReservation(res models.Reservation) (int, error) {

	//best practice for web app to cancel after some timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	sql := `INSERT INTO reservations 
	(first_name, last_name, phone, email, start_date, end_date, 
		room_id, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`

	row := m.DB.QueryRow(ctx, sql,
		res.FirstName,
		res.LastName,
		res.Phone,
		res.Email,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now())

	err := row.Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, err
}

func (m *postgresRepo) AddRoomRestriction(r models.RoomRestriction) error {
	//best practice for web app to cancel after some timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `INSERT INTO room_restrictions 
	(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.Exec(ctx, sql,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now())

	return err
}

// return true means room is available, else room is taken
func (m *postgresRepo) SearchAvailabilityByRoom(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var n int
	sql := `SELECT count(*) FROM room_restrictions 
		WHERE room_id = $1 AND $2 < end_date AND $3 > start_date`

	row := m.DB.QueryRow(ctx, sql, roomID, start, end)
	err := row.Scan(&n)

	return n == 0, err
}

func (m *postgresRepo) SearchAvailabilityAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	sql := `select r.id, r.room_name
	from rooms r
	where r.id not in (
		select rr.room_id from room_restrictions rr
		where '2022-09-21' < rr.end_date and '2022-09-30' > rr.start_date)
	`
	rows, err := m.DB.Query(ctx, sql, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		// if err != nil {
		// 	return rooms, err
		// }
		rooms = append(rooms, room)
	}

	return rooms, err
}

func (m *postgresRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	sql := `SELECT r.id, r.room_name, r.created_at, r.updated_at FROM rooms r WHERE r.room_id = $1`

	row := m.DB.QueryRow(ctx, sql, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		room.UpdatedAt,
	)

	return room, err
}