package godbtdd

import (
	"time"
)

type Blog struct {
	ID        int64
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
}
