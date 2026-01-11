package models

import "time"

type Image struct {
	ImageId   int64
	FileLink  string `binding:"required"`
	FileName  string `binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int64
}
