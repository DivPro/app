package order

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DivPro/app/internal/storage"
	storageerr "github.com/DivPro/app/internal/storage/errors"
	"github.com/DivPro/app/pkg/entity"
)

var (
	ErrInvalidDates    = errors.New("invalid date range")
	ErrRoomUnavailable = errors.New("room is not available for selected dates")
	ErrInternal        = errors.New("internal error")
)

func (s *Service) Create(ctx context.Context, order entity.Order) (*entity.Order, error) {
	var err error
	daysToBook, err := daysBetween(order.From, order.To)
	if err != nil {
		return nil, err
	}

	err = s.storage.Tx(ctx, func(ctx context.Context, tx *storage.DB) error {
		availabilities := make([]*entity.RoomAvailability, 0, len(daysToBook))
		for _, dayToBook := range daysToBook {
			availability, err := tx.FindAvailability(ctx, dayToBook, order.HotelID, order.RoomID)
			if err != nil {
				if errors.Is(err, storageerr.ErrNotFound) {
					continue
				}
				slog.Error("db error", slog.String("error", err.Error()))
			}
			if availability.Quota == 0 {
				continue
			}
			availabilities = append(availabilities, availability)
		}
		if len(daysToBook) != len(availabilities) {
			return ErrRoomUnavailable
		}
		for _, availability := range availabilities {
			availability.Quota--
		}
		if err := tx.CreateOrder(ctx, order); err != nil {
			slog.Error("db error", slog.String("error", err.Error()))
			return ErrInternal
		}
		if err := tx.UpdateAvailabilities(ctx, availabilities); err != nil {
			slog.Error("db error", slog.String("error", err.Error()))
			return ErrInternal
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func daysBetween(from time.Time, to time.Time) ([]time.Time, error) {
	if from.After(to) {
		return nil, ErrInvalidDates
	}

	days := make([]time.Time, 0)
	for d := toDay(from); !d.After(toDay(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days, nil
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}
