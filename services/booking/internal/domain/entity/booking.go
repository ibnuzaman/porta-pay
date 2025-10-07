package entity

import "time"

type BookingStatus string

const (
	StatusCreated   BookingStatus = "CREATED"
	StatusPaid      BookingStatus = "PAID"
	StatusConfirmed BookingStatus = "CONFIRMED"
	StatusExpired   BookingStatus = "EXPIRED"
)

type Booking struct {
	ID         int64         `json:"id" db:"id"`
	UserID     int64         `json:"user_id" db:"user_id"`
	RouteID    int64         `json:"route_id" db:"route_id"`
	Qty        int           `json:"qty" db:"qty"`
	Status     BookingStatus `json:"status" db:"status"`
	PriceTotal int64         `json:"price_total" db:"price_total"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at"`
}
