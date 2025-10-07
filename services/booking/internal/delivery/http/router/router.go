package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/delivery/http/handler"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/delivery/http/middleware"
)

func NewBookingRouter(bookingHandler *handler.BookingHandler) chi.Router {
	r := chi.NewRouter()

	// Apply middleware stack
	for _, mw := range middleware.DefaultStack() {
		r.Use(mw)
	}

	// Health check endpoints
	r.Get("/health", bookingHandler.Health)
	r.Get("/ping", bookingHandler.Health)

	// API endpoints
	r.Route("/api/v1/bookings", func(r chi.Router) {
		r.Post("/", bookingHandler.CreateBooking)
		r.Get("/", bookingHandler.ListBookings)
		r.Get("/{id}", bookingHandler.GetBooking)
		r.Put("/{id}", bookingHandler.UpdateBooking)
		r.Delete("/{id}", bookingHandler.CancelBooking)
	})

	return r
}
