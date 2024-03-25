package storage

import (
	"context"
	"time"

	"github.com/DivPro/app/internal/storage/errors"
	"github.com/DivPro/app/pkg/entity"
)

func (db *DB) FindAvailability(_ context.Context, date time.Time, hotelID string, roomID string) (*entity.RoomAvailability, error) {
	v, ok := db.roomAvailability[createAvailabilityID(hotelID, roomID, date)]
	if !ok {
		return nil, errors.ErrNotFound
	}

	return &entity.RoomAvailability{
		HotelID: hotelID,
		RoomID:  roomID,
		Date:    date,
		Quota:   v.Quota,
	}, nil
}

func (db *DB) CreateOrder(_ context.Context, order entity.Order) error {
	db.orders = append(db.orders, order)

	return nil
}

func (db *DB) UpdateAvailabilities(_ context.Context, availabilities []*entity.RoomAvailability) error {
	for _, v := range availabilities {
		db.roomAvailability[createAvailabilityID(v.HotelID, v.RoomID, v.Date)] = availability{Quota: v.Quota}
	}

	return nil
}
