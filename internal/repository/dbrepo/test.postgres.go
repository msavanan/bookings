package dbrepo

import (
	"errors"
	"time"

	"github.com/msavanan/bookings/internal/models"
)

func (m *postgresTestDBRepo) AllUsers() bool {
	return true
}

func (m *postgresTestDBRepo) InsertReservations(res models.Reservation) (int, error) {

	var newRestrictionId int

	if res.RoomId == 2 {
		return 0, errors.New("wrong restriction ID")
	}

	return newRestrictionId, nil
}

func (m *postgresTestDBRepo) InsertRoomRestrictions(r models.RoomRestriction) error {
	if r.RoomId == 100 {
		return errors.New("wrong restriction ID")
	}
	return nil
}

func (m *postgresTestDBRepo) SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error) {

	return true, nil
}

func (m *postgresTestDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

func (m *postgresTestDBRepo) GetRoomById(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("room id can't be greater than 2")
	}

	return room, nil

}

// Authentication
func (m *postgresTestDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}
func (m *postgresTestDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *postgresTestDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "me@here.ca" {
		return 1, "", nil
	}

	return 0, "", errors.New("invalid email")

}

func (m *postgresTestDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

func (m *postgresTestDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

func (m *postgresTestDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil

}

func (m *postgresTestDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *postgresTestDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *postgresTestDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *postgresTestDBRepo) AllRooms() ([]models.Room, error) {
	return []models.Room{}, nil
}

func (m *postgresTestDBRepo) GetRestrictionForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	var roomRestritions []models.RoomRestriction

	return roomRestritions, nil
}

func (m *postgresTestDBRepo) InsertBlockForRoom(roomId int, startDate time.Time) error {
	return nil
}

func (m *postgresTestDBRepo) DeleteBlockByID(roomId int) error {

	return nil

}
