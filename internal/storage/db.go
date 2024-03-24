package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DivPro/app/pkg/entity"
)

type availabilityID struct {
	HotelID string
	RoomID  string
	Date    time.Time
}

func createAvailabilityID(hotelID, roomID string, date time.Time) availabilityID {
	return availabilityID{
		HotelID: hotelID,
		RoomID:  roomID,
		Date:    date,
	}
}

type availability struct {
	Quota int
}

func (id availabilityID) String() string {
	return fmt.Sprintf("%s:%s:%s", id.HotelID, id.RoomID, id.Date.Format("2006-01-01"))
}

type DB struct {
	mx               sync.Mutex
	orders           []entity.Order
	roomAvailability map[availabilityID]availability
}

func New() *DB {
	return &DB{
		orders: []entity.Order{},
		roomAvailability: map[availabilityID]availability{
			createAvailabilityID("reddison", "lux", date(2024, 1, 1)): {Quota: 1},
			createAvailabilityID("reddison", "lux", date(2024, 1, 2)): {Quota: 1},
			createAvailabilityID("reddison", "lux", date(2024, 1, 3)): {Quota: 1},
			createAvailabilityID("reddison", "lux", date(2024, 1, 4)): {Quota: 1},
			createAvailabilityID("reddison", "lux", date(2024, 1, 5)): {Quota: 0},
		},
	}
}

type TxFn func(ctx context.Context, db *DB) error

func (db *DB) Tx(ctx context.Context, fn TxFn) error {
	db.mx.Lock()
	defer db.mx.Unlock()

	return fn(ctx, db)
}

func date(year, month, day int) time.Time { //nolint:unparam
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
