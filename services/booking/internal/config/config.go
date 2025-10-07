package config

import (
	"github.com/ibnuzaman/porta-pay/pkg/config"
)

// BookingConfig extends the base config with booking-specific settings
type BookingConfig struct {
	*config.Config

	// Booking-specific configurations
	MaxBookingQty      int `env:"MAX_BOOKING_QTY,default=10"`
	BookingExpiryHours int `env:"BOOKING_EXPIRY_HOURS,default=1"`
	AllowCancelHours   int `env:"ALLOW_CANCEL_HOURS,default=2"`
}

// LoadBookingConfig loads booking service configuration
func LoadBookingConfig() (*BookingConfig, error) {
	baseConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	bookingConfig := &BookingConfig{
		Config: &baseConfig,
	}

	return bookingConfig, nil
}
