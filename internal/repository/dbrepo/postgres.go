package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/msavanan/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservations(res models.Reservation) (int, error) {

	var newRestrictionId int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newRestrictionId)

	if err != nil {
		return newRestrictionId, err
	}

	return newRestrictionId, nil
}

func (m *postgresDBRepo) InsertRoomRestrictions(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		r.RestrictionId,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error) {

	var numRows int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT COUNT(id) FROM room_restrictions where  room_id = $1 AND $2 < start_date AND $3 > end_date;
	`

	rows := m.DB.QueryRowContext(ctx, query, roomId, start, end)

	err := rows.Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT r.id, r.room_name
	FROM rooms r
	wHERE r.id NOT IN (SELECT room_id 
		FROM room_restrictions rr 
		WHERE $1 < end_date AND $2 > start_date);
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err

	}

	return rooms, nil

}

func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
	SELECT id, room_name, created_at, updated_at from rooms where id=$1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return room, err
	}

	return room, nil

}

func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err = error(nil)

	query := `SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at 
	FROM users
	WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err = row.Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt)

	return u, err
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	var err = error(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users set id = $1, first_name = $2, last_name = $3, email = $4, access_level = $5, updated_at = $6`

	_, err = m.DB.ExecContext(ctx, query, u.Id, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	return err

}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	var id int
	var hashedPassword string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, password 
	FROM users 
	WHERE email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return id, "", err
	}

	return id, hashedPassword, err
}

func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone,
		 r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		 rm.id, rm.room_name
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		ORDER BY r.start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	err = rows.Err()
	if err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone,
		 r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
		 rm.id, rm.room_name
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		WHERE processed = 0
		ORDER BY r.start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	err = rows.Err()
	if err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone,
		 r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		 rm.id, rm.room_name
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		WHERE r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomId,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil

}

func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	var err = error(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE reservations set first_name = $2, last_name = $3, email = $4, phone = $5, updated_at = $6
	WHERE id = $1
	`

	_, err = m.DB.ExecContext(ctx, query, u.Id, u.FirstName, u.LastName, u.Email, u.Phone, time.Now())

	return err

}

// DeleteReservation deletes reservation by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM reservations WHERE id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE reservations set processed = $1 WHERE id =$2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil

}
