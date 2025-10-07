package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/entity"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/repository"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/service"
)

type bookingUsecase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUsecase(bookingRepo repository.BookingRepository) service.BookingService {
	return &bookingUsecase{
		bookingRepo: bookingRepo,
	}
}

func (uc *bookingUsecase) CreateBooking(ctx context.Context, booking *entity.Booking) error {
	// Business logic validation
	if booking.Qty <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Set default values
	booking.Status = entity.StatusCreated
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	return uc.bookingRepo.Create(ctx, booking)
}

func (uc *bookingUsecase) GetBooking(ctx context.Context, id int64) (*entity.Booking, error) {
	return uc.bookingRepo.GetByID(ctx, id)
}

func (uc *bookingUsecase) UpdateBooking(ctx context.Context, booking *entity.Booking) error {
	// Business logic validation
	existingBooking, err := uc.bookingRepo.GetByID(ctx, booking.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	booking.UpdatedAt = time.Now()
	booking.CreatedAt = existingBooking.CreatedAt // Preserve original creation time

	return uc.bookingRepo.Update(ctx, booking)
}

func (uc *bookingUsecase) CancelBooking(ctx context.Context, id int64) error {
	booking, err := uc.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Business logic: can only cancel if not already confirmed
	if booking.Status == entity.StatusConfirmed {
		return errors.New("cannot cancel confirmed booking")
	}

	booking.Status = entity.StatusExpired
	booking.UpdatedAt = time.Now()

	return uc.bookingRepo.Update(ctx, booking)
}

func (uc *bookingUsecase) ListBookings(ctx context.Context, limit, offset int) ([]*entity.Booking, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return uc.bookingRepo.List(ctx, limit, offset)
}
