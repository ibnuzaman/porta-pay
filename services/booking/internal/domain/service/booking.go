package service

import (
	"context"

	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/entity"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking *entity.Booking) error
	GetBooking(ctx context.Context, id int64) (*entity.Booking, error)
	UpdateBooking(ctx context.Context, booking *entity.Booking) error
	CancelBooking(ctx context.Context, id int64) error
	ListBookings(ctx context.Context, limit, offset int) ([]*entity.Booking, error)
}
