package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/keroda/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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
		&room.UpdatedAt,
	)

	return room, err
}

func (m *postgresRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	sql := `SELECT id, first_name, last_name, email, password, 
		access_level, created_at, updated_at 
		FROM users WHERE id = $1`

	row := m.DB.QueryRow(ctx, sql, id)
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return u, err
}

func (m *postgresRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `UPDATE users SET 
		first_name = $1,
		last_name = $2,
		email = $3,
		access_level = $4,
		updated_at = $5
		WHERE id = $6
		`

	_, err := m.DB.Exec(ctx, sql,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID)

	return err
}

func (m *postgresRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	sql := "SELECT id, apassword FROM users WHERE email = '$1'"

	row := m.DB.QueryRow(ctx, sql, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
		r.created_at, r.updated_at,	rm.id, rm.room_name, r.processed
	FROM reservations
	 left join rooms ON r.room_id = rm.id
	ORDER BY r.start_date ASC`

	rows, err := m.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RoomID,
			&i.Room.RoomName,
			&i.Processed,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
		r.created_at, r.updated_at,	rm.id, rm.room_name
	FROM reservations
	 left join rooms ON r.room_id = rm.id
	 WHERE r.processed = 0
	ORDER BY r.start_date ASC`

	rows, err := m.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RoomID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	sql := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
	r.created_at, r.updated_at,	rm.id, rm.room_name
FROM reservations
 left join rooms ON r.room_id = rm.id
 WHERE r.id = $1`

	row := m.DB.QueryRow(ctx, sql, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.RoomID,
		&res.Room.RoomName,
	)
	// if err != nil {
	// 	return res, err
	// }

	return res, err
}

func (m *postgresRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `UPDATE reservations SET 
		first_name = $1,
		last_name = $2,
		email = $3,
		phone = $4,
		updated_at = $5
		WHERE id = $6
		`

	_, err := m.DB.Exec(ctx, sql,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)

	return err
}

func (m *postgresRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := "DELETE FROM reservations WHERE id = $1"

	_, err := m.DB.Exec(ctx, sql, id)

	return err
}

func (m *postgresRepo) UpdateProcessedRservation(id int, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := "UPDATE reservations SET processed = $1 WHERE id = $2"

	_, err := m.DB.Exec(ctx, sql, processed, id)

	return err
}

func (m *postgresRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	sql := `SELECT id, room_name, created_at, updated_at
	FROM rooms ORDER BY room_name`

	rows, err := m.DB.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Room
		err = rows.Scan(
			&i.ID,
			&i.RoomName,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, i)
	}

	return rooms, err
}

func (m *postgresRepo) GetRoomRestrictionsByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction
	sql := `SELECT id, coalesce(reservation_id, 0), 
		restriction_id ,start_date, end_date
		FROM room_restrictions WHERE $1 < end_date AND $2 >= start_date
			and room_id = $3			
			`
	rows, err := m.DB.Query(ctx, sql, end, start, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.RoomRestriction
		err = rows.Scan(
			&i.ID,
			&i.ReservationID,
			&i.RestrictionID,
			&i.StartDate,
			&i.EndDate,
		)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, i)
	}

	return restrictions, err
}

func (m *postgresRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `
		INSERT INTO room_restrictions 
		(start_date, end_date, room_id, restriction_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,&6)
	`

	_, err := m.DB.Exec(ctx, sql, startDate, startDate.AddDate(0, 0, 1), id, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *postgresRepo) DeleteBlockById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `
		DELETE FROM room_restrictions WHERE id = $1
	`
	_, err := m.DB.Exec(ctx, sql, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
