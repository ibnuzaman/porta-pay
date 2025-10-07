package repository

import (
	"context"

	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/entity"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type postgresBookingRepository struct {
	db *sqlx.DB
}

func NewPostgresBookingRepository(db *sqlx.DB) repository.BookingRepository {
	return &postgresBookingRepository{
		db: db,
	}
}

func (r *postgresBookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	query := `
		INSERT INTO bookings (user_id, route_id, qty, status, price_total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		booking.UserID,
		booking.RouteID,
		booking.Qty,
		booking.Status,
		booking.PriceTotal,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(&booking.ID)

	return err
}

func (r *postgresBookingRepository) GetByID(ctx context.Context, id int64) (*entity.Booking, error) {
	query := `
		SELECT id, user_id, route_id, qty, status, price_total, created_at, updated_at
		FROM bookings
		WHERE id = $1`

	booking := &entity.Booking{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.RouteID,
		&booking.Qty,
		&booking.Status,
		&booking.PriceTotal,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *postgresBookingRepository) Update(ctx context.Context, booking *entity.Booking) error {
	query := `
		UPDATE bookings
		SET user_id = $2, route_id = $3, qty = $4, status = $5, price_total = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		booking.ID,
		booking.UserID,
		booking.RouteID,
		booking.Qty,
		booking.Status,
		booking.PriceTotal,
		booking.UpdatedAt,
	)

	return err
}

func (r *postgresBookingRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresBookingRepository) List(ctx context.Context, limit, offset int) ([]*entity.Booking, error) {
	query := `
		SELECT id, user_id, route_id, qty, status, price_total, created_at, updated_at
		FROM bookings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		booking := &entity.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.RouteID,
			&booking.Qty,
			&booking.Status,
			&booking.PriceTotal,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
