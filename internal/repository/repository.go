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
}
