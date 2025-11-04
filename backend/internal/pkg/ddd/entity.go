package ddd

import "time"

type Entity[ID any] struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Version   int
}
