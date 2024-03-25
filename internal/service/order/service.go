package order

import (
	"context"
	"time"

	"github.com/DivPro/app/internal/storage"
	"github.com/DivPro/app/pkg/entity"
)

type TX interface {
	Tx(ctx context.Context, fn storage.TxFn) error
}

type Order interface {
	FindAvailability(ctx context.Context, date time.Time, hotelID string, roomID string) (*entity.RoomAvailability, error)
	CreateOrder(ctx context.Context, order entity.Order) error
	UpdateAvailabilities(ctx context.Context, availabilities []*entity.RoomAvailability) error
}

//go:generate mockery --name Storage
type Storage interface {
	TX
	Order
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}
