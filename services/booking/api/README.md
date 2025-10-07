# Booking Service API

## Overview
This is the Booking Service API documentation for the Porta Pay microservice.

## Endpoints

### Health Check
- **GET** `/health` - Health check endpoint
- **GET** `/ping` - Alternative health check endpoint

### Bookings
- **POST** `/api/v1/bookings` - Create a new booking
- **GET** `/api/v1/bookings` - List bookings with pagination
- **GET** `/api/v1/bookings/{id}` - Get booking by ID
- **PUT** `/api/v1/bookings/{id}` - Update booking
- **DELETE** `/api/v1/bookings/{id}` - Cancel booking

## Request/Response Examples

### Create Booking
```bash
POST /api/v1/bookings
Content-Type: application/json

{
  "user_id": 123,
  "route_id": 456,
  "qty": 2,
  "price_total": 50000
}
```

### Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 123,
    "route_id": 456,
    "qty": 2,
    "status": "CREATED",
    "price_total": 50000,
    "created_at": "2025-10-07T23:00:00Z",
    "updated_at": "2025-10-07T23:00:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid JSON"
  }
}
```