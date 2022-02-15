package orders

import (
	"time"
)

var UnprocessedStatuses = []string{"NEW", "PROCESSING"}

// Order represents order transfer object.
type Order struct {
	ID         int64
	Number     string
	Status     string
	Accrual    float64
	UserID     int64
	UploadedAt time.Time
}
