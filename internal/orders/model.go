package orders

import (
	"time"
)

type OrderStatus string

var UnprocessedStatuses = []OrderStatus{"NEW", "PROCESSING"}

// Order represents order transfer object.
type Order struct {
	ID         int64
	Number     string
	Status     OrderStatus
	Accrual    float64
	UserID     int64
	UploadedAt time.Time
}
