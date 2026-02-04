package models

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/my-ablum/database"
)

type File struct {
	FileId     int64     `json:"file_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	StorageKey string    `json:"storage_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  int64     `json:"created_by"`
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

func (f *File) Delete() error {
	query := `DELETE FROM files
				WHERE id = $1`

	_, err := database.DB.Exec(query, f.FileId)

	return err
}

func GetFileById(fileId int64) (*File, error) {
	query := `SELECT id, file_name, file_size, mime_type, storage_key, created_at, updated_at, created_by
        		FROM files
       			WHERE id = $1
			`

	row := database.DB.QueryRow(query, fileId)

	var f File

	err := row.Scan(
		&f.FileId,
		&f.FileName,
		&f.FileSize,
		&f.MimeType,
		&f.StorageKey,
		&f.CreatedAt,
		&f.UpdatedAt,
		&f.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &f, err

}

func GetFilesByUserId(userId, cursor, limit int64) ([]File, int64, error) {

	if limit <= 0 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	var rows *sql.Rows
	var err error

	// We assume sorting by Newest First (DESC).
	// If cursor is 0 (first page), just get the top items.
	// If cursor > 0, get items OLDER (smaller ID) than the cursor.

	baseQuery := `
        SELECT id, file_name, file_size, mime_type, storage_key, created_at, updated_at, created_by
        FROM files
        WHERE created_by = $1
    `

	if cursor == 0 {
		query := baseQuery + ` ORDER BY id DESC LIMIT $2`
		rows, err = database.DB.Query(query, userId, limit)

	} else {
		query := baseQuery + ` AND id < $2 ORDER BY id DESC LIMIT $3`
		rows, err = database.DB.Query(query, userId, cursor, limit)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	files := make([]File, 0, limit) //length of slice is 0 but reserve a memory of limit in the bg (for easier appends)

	for rows.Next() {
		var f File
		err := rows.Scan(
			&f.FileId,
			&f.FileName,
			&f.FileSize,
			&f.MimeType,
			&f.StorageKey,
			&f.CreatedAt,
			&f.UpdatedAt,
			&f.CreatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var nextCursor int64 = 0

	if len(files) > 0 {
		nextCursor = files[len(files)-1].FileId
	}

	return files, nextCursor, nil
}
