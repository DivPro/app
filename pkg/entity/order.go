package entity

import "time"

type Order struct {
	HotelID   string    `json:"hotel_id" validate:"required"`
	RoomID    string    `json:"room_id" validate:"required"`
	UserEmail string    `json:"email" validate:"required"`
	From      time.Time `json:"from" validate:"required"`
	To        time.Time `json:"to" validate:"required"`
}
