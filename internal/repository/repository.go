package repository

import (
	"time"

	"github.com/msavanan/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservations(res models.Reservation) (int, error)

	InsertRoomRestrictions(r models.RoomRestriction) error

	SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)

	//Authentication
	GetUserById(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)
	UpdateUser(u models.User) error

	//Reservations
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)

	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
