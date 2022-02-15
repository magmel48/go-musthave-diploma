package orders

import (
	"time"
)

var UnprocessedStatuses = []string{"NEW", "PROCESSING"}

// Order represents order transfer object.
type Order struct {
	ID         int64     `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UserID     int64     `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}
