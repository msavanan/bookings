package dbrepo

import (
	"context"
	"time"

	"github.com/msavanan/bookings/internal/models"
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
