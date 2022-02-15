package orders

import (
	"time"
)

type OrderStatus string

// Order represents order transfer object.
type Order struct {
	ID         int64
	Number     string
	Status     OrderStatus
	Accrual    float64
	UserID     int64
	UploadedAt time.Time
}
