package models

import (
	"time"

	"example.com/my-ablum/database"
)

type Image struct {
	ImageId   int64
	FileLink  string    `json:"file_link" binding:"required"`
	FileName  string    `json:"file_name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
}

func (f *Image) Save() error {

	query := `INSERT INTO files (file_name, file_link, created_by)
			 VALUES($1, $2, $3)
			 RETURNING id;
			`
	row := database.DB.QueryRow(query, f.FileName, f.FileLink, f.CreatedBy)
	err := row.Scan(&f.ImageId)

	return err

}
