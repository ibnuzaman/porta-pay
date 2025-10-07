package repository

import (
	"context"

	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/entity"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByID(ctx context.Context, id int64) (*entity.Booking, error)
	Update(ctx context.Context, booking *entity.Booking) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.Booking, error)
}
