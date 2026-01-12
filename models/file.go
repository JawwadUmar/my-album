package models

import (
	"time"

	"example.com/my-ablum/database"
)

type File struct {
	FileId     int64
	FileName   string
	FileSize   int64
	MimeType   string
	StorageKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CreatedBy  int64
}

func (f *File) Save() error {

	query := `INSERT INTO files (file_name, file_size, mime_type, storage_key, created_by)
			 VALUES($1, $2, $3, $4, $5)
			 RETURNING id;
			`
	row := database.DB.QueryRow(query, f.FileName, f.FileSize, f.MimeType, f.StorageKey, f.CreatedBy)
	err := row.Scan(&f.FileId)

	return err

}
