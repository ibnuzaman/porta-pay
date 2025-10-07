package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ibnuzaman/porta-pay/pkg/response"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/entity"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/domain/service"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking entity.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		response.BadRequest(w, "Invalid JSON")
		return
	}

	if err := h.bookingService.CreateBooking(r.Context(), &booking); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, booking)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	booking, err := h.bookingService.GetBooking(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Booking not found")
		return
	}

	response.Success(w, http.StatusOK, booking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	var booking entity.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		response.BadRequest(w, "Invalid JSON")
		return
	}

	booking.ID = id
	if err := h.bookingService.UpdateBooking(r.Context(), &booking); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, http.StatusOK, booking)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	if err := h.bookingService.CancelBooking(r.Context(), id); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Booking cancelled successfully"})
}

func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	bookings, err := h.bookingService.ListBookings(r.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, http.StatusOK, bookings)
}

func (h *BookingHandler) Health(w http.ResponseWriter, r *http.Request) {
	healthData := map[string]string{
		"status":  "ok",
		"service": "booking",
	}

	response.Success(w, http.StatusOK, healthData)
}
