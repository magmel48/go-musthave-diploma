package orders

import "time"

// Order represents order transfer object.
type Order struct {
	ID         int64
	Number     string
	Status     string
	UserID     int64
	UploadedAt time.Time
}
